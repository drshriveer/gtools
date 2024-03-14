package set_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/drshriveer/gtools/set"
)

func TestSet_MarshalJSON(t *testing.T) {
	type jsonMarshalType struct {
		StringSet set.Set[string] `json:"string_set"`
	}
	initial := jsonMarshalType{
		StringSet: set.Make("s0", "s1", "s3"),
	}
	binary, err := json.Marshal(initial)
	require.NoError(t, err)
	unmarshal := jsonMarshalType{}
	require.NoError(t, json.Unmarshal(binary, &unmarshal))
	assert.Equal(t, initial, unmarshal)
}

func TestSet_MarshalYAML(t *testing.T) {
	type yamlMarshalType struct {
		StringSet set.Set[string] `yaml:"string_set"`
	}
	initial := yamlMarshalType{
		StringSet: set.Make("s0", "s1", "s3"),
	}
	binary, err := yaml.Marshal(initial)
	require.NoError(t, err)
	unmarshal := yamlMarshalType{}
	require.NoError(t, yaml.Unmarshal(binary, &unmarshal))
	assert.Equal(t, initial, unmarshal)
}

func TestSet_NilHandling(t *testing.T) {
	var s1 set.Set[string]
	assert.Nil(t, s1.Slice())
	assert.False(t, s1.Remove("hi"))
	assert.False(t, s1.RemoveSet(set.Make("bye")))
	assert.False(t, s1.Has("anything"))
	assert.False(t, s1.HasAny("anything"))

	assert.True(t, s1.Add("hi"))
	assert.True(t, s1.Has("hi"))
	assert.False(t, s1.Has("hi", "bye"))
	assert.True(t, s1.HasAny("hi", "bye"))
	assert.True(t, s1.Remove("hi", "bye"))
	assert.False(t, s1.Has("hi"))

	s1 = nil
	assert.True(t, s1.AddSet(set.Make("hi", "bye")))
	assert.True(t, s1.Has("hi"))
	assert.True(t, s1.Has("bye", "hi"))
	assert.True(t, s1.Remove("hi", "cry"))
	assert.False(t, s1.Has("hi"))
	assert.True(t, s1.Has("bye"))
}
