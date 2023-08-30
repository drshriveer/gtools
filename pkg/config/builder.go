package config

import (
	"encoding"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"reflect"

	"github.com/puzpuzpuz/xsync/v2"
	"gopkg.in/yaml.v3"

	"github.com/drshriveer/gcommon/pkg/enum"
	"github.com/drshriveer/gcommon/pkg/errors"
	"github.com/drshriveer/gcommon/pkg/set"
)

const defaultKey = "default"

var ErrFailedParsing errors.Factory = &errors.GError{
	Name:    "ErrFailedParsing",
	Message: "failed to read or parse configuration",
}

type Dimension struct {
	Default  enum.Enum
	FlagName string
	ParseEnv bool

	parsed enum.Enum
}

func (d *Dimension) initFlag() error {
	d.parsed = d.Default
	if !d.ParseEnv {
		return nil
	}

	usage := fmt.Sprintf("%s (default=%s): configuration dimension valid options: %s",
		d.FlagName, d.Default, d.Default.StringValues())

	eType := reflect.TypeOf(d.Default)
	ptrVal, ok := reflect.New(eType).Interface().(encoding.TextUnmarshaler)
	if !ok {
		return ErrFailedParsing.Include(
			"genum %T does not implement encoding.TextUnmarshaler as required",
			d.Default)
	}

	// first look for flags that have already been registered..
	// if so, this is probably a testing environment, so skip the flag registration.
	// long term need to decide if we want to disallow this for safety?
	// FIXME: Gavin! test the behavior here... Do we need to tie into the parse function
	//        in a different way to receive updates if parse is called?
	//        or is this just all nuts, and that's why I was originally approaching this w/o
	//        a builder?
	//        The whole flag part is ... maybe problematic.
	//        Maybe there's an easier way?
	if flag.Lookup(d.FlagName) != nil {
		return nil
	}

	flag.Func(d.FlagName, usage, func(s string) error {
		if err := ptrVal.UnmarshalText([]byte(s)); err != nil {
			return err
		}
		d.parsed, ok = unptr(ptrVal).(enum.Enum)
		if !ok {
			return ErrFailedParsing.WithStack()
		}
		return nil
	})

	return nil
}

func (d *Dimension) get() (enum.Enum, error) {
	if !flag.Parsed() {
		flag.Parse()
	}
	return d.parsed, nil
}

type Builder struct {
	Dimensions []Dimension
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithDimension(name string, defaultVal enum.Enum) *Builder {
	d := Dimension{
		Default:  defaultVal,
		FlagName: name,
		ParseEnv: true,
		parsed:   defaultVal,
	}
	if err := d.initFlag(); err != nil {
		panic(err)
	}
	b.Dimensions = append(b.Dimensions, d)
	return b
}

func (b *Builder) FromFile(fileSystem fs.FS, filename string) (*Config, error) {
	f, err := fileSystem.Open(filename)
	if err != nil {
		return nil, ErrFailedParsing.Convert(err)
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, ErrFailedParsing.Convert(err)
	}

	return b.FromBytes(bytes)
}

func (b *Builder) FromBytes(bytes []byte) (*Config, error) {
	data := make(map[string]any)
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return nil, ErrFailedParsing.Convert(err)
	}

	d, err := reduceAny(data, b.Dimensions, 0)
	if err != nil {
		return nil, err
	}
	result, ok := d.(map[string]any)
	if !ok {
		return nil, ErrFailedParsing.Include("unexpected non-map result")
	}
	cfg := &Config{
		cached: xsync.NewMapOf[any](),
		data:   result,
	}
	return cfg, nil
}

// XXX: I think last time I reduced the keys I traced them down to a bottom value
// to see if they successfully resolved.
// . Somehow that doesn't seem right/efficient. Why should I traverse unnecessary options?
// . But maybe that was a form of validation that the structure I'm tracing is not,
//   for example, a map type of a dimension type. The best example I can think of for this use case
//   is regionalization. Sometimes you want a region to be a dimension (setting x is enabled in region y)
//   but sometimes you want to separate something by a _discoverable_ region (regional client addresses)
//  . a different version of safety could use special characters to difference a enum from map key.
//    . the parser/builder step could convert the differenced enums to parsable characters.
//      . then

func reduceAny(in any, dimensions []Dimension, dIndex int) (any, error) {
	switch v := in.(type) {
	case map[string]any:
		for i := dIndex; i < len(dimensions); i++ {
			r, err := reduce(v, dimensions, i)
			if err != nil || !reflect.DeepEqual(r, v) {
				return r, err
			}
		}
	case []any:
		for i, el := range v {
			var err error
			v[i], err = reduceAny(el, dimensions, dIndex)
			if err != nil {
				return nil, err
			}
		}
		// I'm returning an ANY here to do the map reduction in place
		// but this is conflicting with the non-redusable case.
		// return reduce(v, dimensions, dIndex)
	}
	return in, nil
}

func reduce(in map[string]any, dimensions []Dimension, dIndex int) (any, error) {
	if dIndex+1 > len(dimensions) {
		return in, nil
	}
	dimension := dimensions[dIndex]
	// check if this a valid dimension to reduce.
	// if it is, grab the correct one and reduce the rest.
	keys, hasDefault := keySet(in)
	keys.Remove(dimension.Default.StringValues()...)
	if len(keys) != 0 {
		for k, v := range in {
			var err error
			in[k], err = reduceAny(v, dimensions, dIndex)
			if err != nil {
				return nil, err
			}
		}
		// NOT reducable with this dimension. need to try next,
		return in, nil
	}
	// otherwise this is reducable.
	// case 1: we have the dimension's key. Simply follow it.
	if v, ok := in[dimension.parsed.String()]; ok {
		return reduceAny(v, dimensions, dIndex+1)
	}
	// case 2: we have default
	if hasDefault {
		return reduceAny(in[defaultKey], dimensions, dIndex+1)
	}

	// case 3: we have no default, and no match...
	// There are sort of two options here.
	// 1. This is just a completely invalid config
	// 2. These keys are meant to be part of a map... i.e. intentionally missing properties.
	// ...going with #1.
	keys, hasDefault = keySet(in)
	return nil, ErrFailedParsing.Include(
		"broken dimension key! %T dimensions identified around keys %s, but no `default` or `%s` value found.",
		dimension.Default, keys.Slice(), dimension.parsed)
}

func keySet(in map[string]any) (set.Set[string], bool) {
	hasDefault := false
	result := make(set.Set[string], len(in))
	for k := range in {
		if k == defaultKey {
			hasDefault = true
		} else {
			result.Add(k)
		}
	}
	return result, hasDefault
}

func unptr(in any) any {
	v := reflect.ValueOf(in)
	switch v.Kind() {
	case reflect.Pointer:
		return v.Elem().Interface()
	default:
		return v.Interface()
	}
}
