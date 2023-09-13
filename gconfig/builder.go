package gconfig

import (
	"encoding"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/puzpuzpuz/xsync/v2"

	"github.com/drshriveer/gtools/genum"
	"github.com/drshriveer/gtools/gerrors"
	"github.com/drshriveer/gtools/rutils"
	"github.com/drshriveer/gtools/set"
)

const defaultKey = "default"

// ErrFailedParsing is returned if there are errors parsing a config file.
var ErrFailedParsing gerrors.Factory = &gerrors.GError{
	Name:    "ErrFailedParsing",
	Message: "failed to read or parse configuration",
}

// A dimension is a running parameter that instructs which configuration to use.
type dimension struct {
	// defaultVal is the default enum value if none is supplied.
	// This is also used to determine the
	defaultVal genum.Enum

	// flagName is the name of the environment flag to parse when parseEnv is true.
	flagName string

	// parseEnv, if true, will parse the dimension as an environment flag.
	parseEnv bool

	parsed genum.Enum
}

func (d *dimension) initFlag() error {
	d.parsed = d.defaultVal
	if !d.parseEnv {
		return nil
	}

	usage := fmt.Sprintf("%s (default=%s): configuration dimension valid options: %s",
		d.flagName, d.defaultVal, d.defaultVal.StringValues())

	eType := reflect.TypeOf(d.defaultVal)
	ptrVal, ok := reflect.New(eType).Interface().(encoding.TextUnmarshaler)
	if !ok {
		return ErrFailedParsing.Include(
			"genum %T does not implement encoding.TextUnmarshaler as required",
			d.defaultVal)
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
	if flag.Lookup(d.flagName) != nil {
		return nil
	}

	flag.Func(d.flagName, usage, func(s string) error {
		if err := ptrVal.UnmarshalText([]byte(s)); err != nil {
			return err
		}
		d.parsed, ok = rutils.Unptr(ptrVal).(genum.Enum)
		if !ok {
			return ErrFailedParsing.WithStack()
		}
		return nil
	})

	return nil
}

func (d *dimension) get() genum.Enum {
	if !flag.Parsed() {
		flag.Parse()
	}
	return d.parsed
}

// Builder is a configuration builder.
type Builder struct {
	// An ordered set of dimensions to switch a configuration on.
	dimensions []dimension
}

// NewBuilder returns a new builder instance.
func NewBuilder() *Builder {
	return &Builder{}
}

// WithDimension adds a new dimension to switch configurations on. By default `parseEnv` will be true when using this method.
func (b *Builder) WithDimension(name string, defaultVal genum.Enum) *Builder {
	d := dimension{
		defaultVal: defaultVal,
		flagName:   name,
		parseEnv:   true,
		parsed:     defaultVal,
	}
	if err := d.initFlag(); err != nil {
		panic(err)
	}
	b.dimensions = append(b.dimensions, d)
	return b
}

// FromFile takes a file system and a path to a configuration file to parse a Config from.
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

// FromBytes takes configuration file bytes and parses a Config object from them.
func (b *Builder) FromBytes(bytes []byte) (*Config, error) {
	data := make(map[string]any)
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return nil, ErrFailedParsing.Convert(err)
	}

	d, err := reduceAny(data, b.dimensions, 0)
	if err != nil {
		return nil, err
	}
	result, ok := d.(map[string]any)
	if !ok {
		return nil, ErrFailedParsing.Include("unexpected non-map result")
	}

	dims := make(map[reflect.Type]genum.Enum, len(b.dimensions))
	for _, d := range b.dimensions {
		dims[reflect.TypeOf(d.defaultVal)] = d.get()
	}

	cfg := &Config{
		dimensions: dims,
		cached:     xsync.NewMapOf[any](),
		data:       result,
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

func reduceAny(in any, dimensions []dimension, dIndex int) (any, error) {
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

func reduce(in map[string]any, dimensions []dimension, dIndex int) (any, error) {
	if dIndex+1 > len(dimensions) {
		return in, nil
	}
	dimension := dimensions[dIndex]
	// check if this a valid dimension to reduce.
	// if it is, grab the correct one and reduce the rest.
	keys, hasDefault := keySet(in)
	keys.Remove(dimension.defaultVal.StringValues()...)
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
	keys, _ = keySet(in)
	return nil, ErrFailedParsing.Include(
		"broken dimension key! %T dimensions identified around keys %s, but no `default` or `%s` value found.",
		dimension.defaultVal, keys.Slice(), dimension.parsed)
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