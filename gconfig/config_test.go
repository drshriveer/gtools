package gconfig_test

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gtools/gconfig/internal"

	"github.com/drshriveer/gtools/gconfig"
)

//go:embed internal/*.yaml
var testFS embed.FS

type testStruct struct {
	Pi       float64       `yaml:"pi"`
	E        float64       `yaml:"e"`
	Duration time.Duration `yaml:"duration"`
	UseReal  bool          `yaml:"useReal"`
	Name     string        `yaml:"name"`
}

func TestFlagParsing(t *testing.T) {
	builder := gconfig.NewBuilder().
		WithDimension("d1", internal.D1a).
		WithDimension("d2", internal.D2a)

	require.NoError(t, flag.Set("d1", internal.D1c.String()))
	require.NoError(t, flag.Set("d2", internal.D2e.String()))
	cfg, err := builder.FromFile(testFS, "internal/test.yaml")
	require.NoError(t, err)
	assert.Equal(t, internal.D1c, gconfig.GetDimension[internal.DimensionOne](cfg))
	assert.Equal(t, internal.D2e, gconfig.GetDimension[internal.DimensionTwo](cfg))
}

func TestEnvParsing(t *testing.T) {
	require.NoError(t, os.Setenv("D1", internal.D1c.String()))
	require.NoError(t, os.Setenv("d2", internal.D2e.String()))
	t.Cleanup(func() {
		require.NoError(t, os.Unsetenv("D1"))
		require.NoError(t, os.Unsetenv("d2"))
	})
	cfg, err := gconfig.NewBuilder().
		WithDimension("d1", internal.D1a).
		WithDimension("d2", internal.D2a).
		FromFile(testFS, "internal/test_template_env_vars.yaml")

	require.NoError(t, err)
	assert.Equal(t, internal.D1c, gconfig.GetDimension[internal.DimensionOne](cfg))
	assert.Equal(t, internal.D2e, gconfig.GetDimension[internal.DimensionTwo](cfg))

	// TEST env var parsing from templates defined in the config.yaml.
	assert.Equal(t, internal.D1c.String(), gconfig.MustGet[string](cfg, "envTemplating.firstDim"))
	assert.Equal(t, internal.D2e.String(), gconfig.MustGet[string](cfg, "envTemplating.secondDim"))
	assert.Equal(t, "default value", gconfig.MustGet[string](cfg, "envTemplating.switchedDim"))
	assert.Equal(t, internal.D1c, gconfig.MustGet[internal.DimensionOne](cfg, "envTemplating.firstDim"))
	assert.Equal(t, internal.D2e, gconfig.MustGet[internal.DimensionTwo](cfg, "envTemplating.secondDim"))
}

func TestDimensions(t *testing.T) {
	standardStructDefault := testStruct{
		Pi:       3,
		E:        2,
		Duration: 2 * time.Minute,
		UseReal:  false,
		Name:     "default name",
	}
	standardStructD1A := testStruct{
		Pi:       3.14159,
		E:        2.71828,
		Duration: 3 * time.Minute,
		UseReal:  true,
		Name:     "a name",
	}
	nestedDimStructD1A := testStruct{
		Pi:       3.14159,
		E:        2.71828,
		Duration: 2*time.Second + time.Millisecond,
		UseReal:  true,
		Name:     "D1a name",
	}
	nestedDimStructD1AD2B := testStruct{
		Pi:       3.14159,
		E:        2.71828,
		Duration: 2*time.Second + time.Millisecond,
		UseReal:  false,
		Name:     "D1a name",
	}
	nestedDimStructD1B := testStruct{
		Pi:       3.14159,
		E:        2.71828,
		Duration: 2*time.Second + time.Millisecond,
		UseReal:  true,
		Name:     "D1b name",
	}
	nestedDimStructD1DD2B := testStruct{
		Pi:       3.14159,
		E:        2.71828,
		Duration: 2*time.Second + time.Millisecond,
		UseReal:  false,
		Name:     "default",
	}

	tests := []struct {
		d1       internal.DimensionOne
		d2       internal.DimensionTwo
		cfgTests map[string]tester
	}{
		{
			d1: internal.D1a,
			d2: internal.D2a,
			cfgTests: map[string]tester{
				"dimensions.valid.v1":    testGetter[string]("v1:D1a.D2a", ""),
				"dimensions.valid.v2":    testGetter[string]("v2:default", ""),
				"dimensions.valid.v3":    testGetter[[]string]([]string{"item.1", "item.2", "item.D1a"}, []string{}),
				"struct.standard":        testGetter[testStruct](standardStructD1A, testStruct{}),
				"struct.nestedDimension": testGetter[testStruct](nestedDimStructD1A, testStruct{}),
			},
		},
		{
			d1: internal.D1a,
			d2: internal.D2b,
			cfgTests: map[string]tester{
				"dimensions.valid.v1":    testGetter[string]("v1:D1a.default", ""),
				"dimensions.valid.v2":    testGetter[string]("v2:D2b", ""),
				"dimensions.valid.v3":    testGetter[[]string]([]string{"item.1", "item.2", "item.D1a"}, []string{}),
				"struct.standard":        testGetter[testStruct](standardStructD1A, testStruct{}),
				"struct.nestedDimension": testGetter[testStruct](nestedDimStructD1AD2B, testStruct{}),
			},
		},
		{
			d1: internal.D1b,
			d2: internal.D2a,
			cfgTests: map[string]tester{
				"dimensions.valid.v1":    testGetter[string]("v1:D1b", ""),
				"dimensions.valid.v2":    testGetter[string]("v2:default", ""),
				"dimensions.valid.v3":    testGetter[[]string]([]string{"item.1", "item.2", "item.D1b"}, []string{}),
				"struct.standard":        testGetter[testStruct](standardStructDefault, testStruct{}),
				"struct.nestedDimension": testGetter[testStruct](nestedDimStructD1B, testStruct{}),
			},
		},
		{
			d1: internal.D1d,
			d2: internal.D2b,
			cfgTests: map[string]tester{
				"dimensions.valid.v1": testGetter[string]("v1:default", ""),
				"dimensions.valid.v2": testGetter[string]("v2:D2b", ""),
				"dimensions.valid.v3": testGetter[[]string]([]string{"item.1", "item.2", "item.default"},
					[]string{}),
				"struct.standard":        testGetter[testStruct](standardStructDefault, testStruct{}),
				"struct.nestedDimension": testGetter[testStruct](nestedDimStructD1DD2B, testStruct{}),
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("cfg: d1=%s, d2=%s", test.d1, test.d2), func(t *testing.T) {
			cfg, err := gconfig.NewBuilder().
				WithDimension("d1", test.d1).
				WithDimension("d2", test.d2).
				FromFile(testFS, "internal/test.yaml")
			require.NoError(t, err)

			t.Run("dimensions", func(t *testing.T) {
				require.NotPanics(t, func() {
					v1 := gconfig.GetDimension[internal.DimensionOne](cfg)
					assert.Equal(t, test.d1, v1)
					v2 := gconfig.GetDimension[internal.DimensionTwo](cfg)
					assert.Equal(t, test.d2, v2)
				})
				assert.Panics(t, func() {
					gconfig.GetDimension[internal.DimensionThree](cfg)
				})
			})

			for k, tst := range test.cfgTests {
				t.Run(k, func(t *testing.T) {
					tst(t, cfg, k)
				})
			}
		})
	}
}

type tester func(t *testing.T, cfg *gconfig.Config, key string)

func TestGetters(t *testing.T) {
	cfg, err := gconfig.NewBuilder().
		WithDimension("d1", internal.D1b).
		WithDimension("d2", internal.D2c).
		FromFile(testFS, "internal/test.yaml")
	require.NoError(t, err)

	tests := []struct {
		key     string
		testers []tester
	}{
		{
			key: "scalars.numbers.float.invalid",
			testers: []tester{
				testGetter[float64](1.11, 1.11),
				testGetter[float32](1.11, 1.11),
				testGetter[*float64](ptr(1.11), ptr(1.11)),
			},
		},
		{
			key: "scalars.numbers.float",
			testers: []tester{
				testGetter[float64](3.14159, 1.11),
				testGetter[float32](3.14159, 1.11),
				testGetter[*float64](ptr(3.14159), ptr(1.11)),
			},
		},
		{
			key: "scalars.numbers.int.invalid",
			testers: []tester{
				testGetter[int](22, 22),
				testGetter[*int](ptr(22), ptr(22)),
			},
		},
		{
			key: "scalars.numbers.int",
			testers: []tester{
				testGetter[int](271, 1),
				testGetter[int64](271, 1),
				testGetter[uint](271, 1),
				testGetter[uint64](271, 1),
				testGetter[*uint64](ptr[uint64](271), ptr[uint64](1)),
			},
		},
		{
			key: "scalars.strings.invalid",
			testers: []tester{
				testGetter[string]("invalid value!", "invalid value!"),
				testGetter[*string](ptr("invalid value!"), ptr("invalid value!")),
			},
		},
		{
			key: "scalars.strings",
			testers: []tester{
				testGetter[string]("a.value.here", "invalid value!"),
				testGetter[*string](ptr("a.value.here"), ptr("invalid value!")),
			},
		},
		{
			key: "scalars.bools.invalid",
			testers: []tester{
				testGetter[bool](false, false),
				testGetter[*bool](ptr(false), ptr(false)),
			},
		},
		{
			key: "scalars.bools",
			testers: []tester{
				testGetter[*bool](ptr(true), ptr(false)),
				testGetter[bool](true, false),
			},
		},
		// slices:
		{
			key: "slices.numbers.float.invalid",
			testers: []tester{
				testGetter[[]float64]([]float64{1.11}, []float64{1.11}),
			},
		},
		{
			key: "slices.numbers.float",
			testers: []tester{
				testGetter[[]float64]([]float64{3.14159, 2.71}, []float64{1.11}),
				testGetter[[]float32]([]float32{3.14159, 2.71}, []float32{1.11}),
			},
		},
		{
			key: "slices.numbers.int.invalid",
			testers: []tester{
				testGetter[[]int]([]int{22}, []int{22}),
			},
		},
		{
			key: "slices.numbers.int",
			testers: []tester{
				testGetter[[]int]([]int{3173, 271, 365, 1}, []int{22}),
				testGetter[[]int64]([]int64{3173, 271, 365, 1}, []int64{22}),
				testGetter[[]uint]([]uint{3173, 271, 365, 1}, []uint{22}),
				testGetter[[]uint64]([]uint64{3173, 271, 365, 1}, []uint64{22}),
			},
		},
		{
			key: "slices.strings.invalid",
			testers: []tester{
				testGetter[[]string]([]string{"invalid value!"}, []string{"invalid value!"}),
			},
		},
		{
			key: "slices.strings",
			testers: []tester{
				testGetter[[]string](
					[]string{"a.value.here", "b.value.here", "c.value.here"},
					[]string{"invalid value!"},
				),
				testGetter[[]*string](
					[]*string{ptr("a.value.here"), ptr("b.value.here"), ptr("c.value.here")},
					[]*string{ptr("invalid value!")},
				),
			},
		},
		{
			key: "slices.bools.invalid",
			testers: []tester{
				testGetter[[]bool]([]bool{false}, []bool{false}),
			},
		},
		{
			key: "slices.bools",
			testers: []tester{
				testGetter[[]bool]([]bool{true, true, false, false}, []bool{false}),
			},
		},
		// maps:
		{
			key: "maps.numbers.invalid",
			testers: []tester{
				testGetter[map[int]string](map[int]string{42: "the answer"}, map[int]string{42: "the answer"}),
			},
		},
		{
			key: "maps.numbers",
			testers: []tester{
				testGetter[map[int]string](map[int]string{1: "thing1", 2: "thing2"}, map[int]string{}),
				testGetter[map[int32]*string](
					map[int32]*string{1: ptr("thing1"), 2: ptr("thing2")},
					map[int32]*string{},
				),
			},
		},
		{
			key: "maps.strings.invalid",
			testers: []tester{
				testGetter[map[string]string](
					map[string]string{"hello": "goodbye"},
					map[string]string{"hello": "goodbye"},
				),
			},
		},
		{
			key: "maps.strings",
			testers: []tester{
				testGetter[map[string]string](
					map[string]string{"s1": "true", "s2": "false", "s3": "false"},
					map[string]string{},
				),
				testGetter[map[string]bool](
					map[string]bool{"s1": true, "s2": false, "s3": false},
					map[string]bool{},
				),
			},
		},
		{
			key: "maps.bools.invalid",
			testers: []tester{
				testGetter[map[bool]string](map[bool]string{false: "f"}, map[bool]string{false: "f"}),
			},
		},
		{
			key: "maps.bools",
			testers: []tester{
				testGetter[map[bool]string](map[bool]string{true: "1", false: "0"}, map[bool]string{}),
				testGetter[map[bool]int](map[bool]int{true: 1, false: 0}, map[bool]int{}),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			for _, tst := range test.testers {
				tst(t, cfg, test.key)
			}
		})
	}
}

func testGetter[T any](expected T, defaultV T) tester {
	return func(t *testing.T, cfg *gconfig.Config, key string) {
		t.Run(fmt.Sprintf("Get[%T]", *new(T)), func(t *testing.T) {
			result, err := gconfig.Get[T](cfg, key)
			if assert.ObjectsAreEqual(expected, defaultV) {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expected, result)
			}
		})
		t.Run(fmt.Sprintf("MustGet[%T]", *new(T)), func(t *testing.T) {
			fn := func() {
				assert.Equal(t, expected, gconfig.MustGet[T](cfg, key))
			}
			if assert.ObjectsAreEqual(expected, defaultV) {
				assert.Panics(t, fn)
			} else {
				assert.NotPanics(t, fn)
			}
		})
		t.Run(fmt.Sprintf("GetOrDefault[%T]", *new(T)), func(t *testing.T) {
			result := gconfig.GetOrDefault[T](cfg, key, defaultV)
			assert.Equal(t, expected, result)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
