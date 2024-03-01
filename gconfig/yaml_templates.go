package gconfig

import (
	"os"
	"regexp"
)

type templateVariable interface {
	MatchAndResolve(in string) (string, bool, error)
}

var templates = []templateVariable{
	&envVarTmpl{},
}

type envVarTmpl struct{}

var envVarTmplMatcher = regexp.MustCompile(`^\$\{\{\s*env:\s*(\w+)\s*\}\}$`)

func (envVarTmpl) MatchAndResolve(in string) (out string, ok bool, err error) {
	out = in
	matches := envVarTmplMatcher.FindStringSubmatch(in)
	if len(matches) == 0 {
		return out, false, nil
	}
	envVarName := matches[1]
	out, ok = os.LookupEnv(envVarName)
	if !ok {
		return out, false, ErrFailedParsing.Msg(
			"templated environment environment variable %s not found in env",
			envVarName)
	}
	return out, true, nil
}

func parseTemplatedElements[T any](in T) (out T, err error) {
	switch v := any(in).(type) {
	case string:
		for _, template := range templates {
			temp, ok, err := template.MatchAndResolve(v)
			if err != nil {
				return out, err
			} else if ok {
				return any(temp).(T), nil
			}
			// else try next template.
		}
	case map[string]any:
		for k, el := range v {
			v[k], err = parseTemplatedElements(el)
			if err != nil {
				return out, err
			}
		}
	case []any:
		for i, el := range v {
			v[i], err = parseTemplatedElements(el)
			if err != nil {
				return out, err
			}
		}
	}
	return in, nil
}
