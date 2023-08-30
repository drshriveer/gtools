package config

import (
	"fmt"
	"strings"

	"github.com/puzpuzpuz/xsync/v2"
	"gopkg.in/yaml.v3"

	"github.com/drshriveer/gcommon/pkg/errors"
)

var ErrConfigFailure errors.Factory = &errors.GError{
	Name:    "ErrConfigFailure",
	Message: "failed to read value",
}

type Config struct {
	cached *xsync.MapOf[string, any]
	data   map[string]any
}

// Get fetches a value from the config and returns an error if there is a problem.
func Get[T any](cfg *Config, key string) (T, error) {
	return getFromCache[T](cfg, key)
}

// MustGet fetches a value from the config and panics if there are any issues.
func MustGet[T any](cfg *Config, key string) T {
	v, err := getFromCache[T](cfg, key)
	if err != nil {
		panic(err)
	}
	return v
}

// GetOrDefault fetches a value from the config and uses the default value if no value was
// found.
func GetOrDefault[T any](cfg *Config, key string, defaultV T) T {
	v, err := getFromCache[T](cfg, key)
	if err != nil {
		return defaultV
	}
	return v
}

func getFromCache[T any](cfg *Config, key string) (T, error) {
	var err error
	var r T
	k := key + fmt.Sprintf("%T", r) // add type to key to prevent complicated conversions.
	v, _ := cfg.cached.Compute(k, func(oldValue any, loaded bool) (newValue any, delete bool) {
		if loaded {
			return oldValue, false
		}
		oldValue, err = extractAndConvert[T](cfg.data, key)
		if err != nil {
			return oldValue, true
		}
		return oldValue, false
	})

	if err != nil {
		return r, err
	}

	return v.(T), nil
}

func extractAndConvert[T any](m map[string]any, key string) (T, error) {
	// TODO: check env overrides.

	paths := strings.Split(key, ".")
	result := *new(T)
	v, ok := extract(m, paths)
	if !ok {
		return result, ErrConfigFailure.Include("key `%s` not found", key)
	}

	bytes, err := yaml.Marshal(v)
	if err != nil {
		return result, ErrConfigFailure.Include("key `%s` failed conversion back to yaml %+v", key, err)
	}

	err = yaml.Unmarshal(bytes, &result)
	if err != nil {
		return result, ErrConfigFailure.Convert(err)
	}
	return result, nil
}

func extract(m map[string]any, keys []string) (any, bool) {
	var last any = nil
	ok := false
	for i, k := range keys {
		last, ok = m[k]
		if !ok {
			return nil, false
		}
		mOK := false
		m, mOK = last.(map[string]any)
		if !mOK && i < len(keys)-1 {
			return nil, false
		}
	}

	return last, ok
}
