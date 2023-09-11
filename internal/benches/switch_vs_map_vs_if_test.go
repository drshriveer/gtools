package benches_test

import (
	"fmt"
	"testing"
)

const (
	v1, v1s = 1, "v1"
	v2, v2s = 2, "v2"
	v3, v3s = 3, "v3"
	v4, v4s = 4, "v4"

	lookupI, lookupS = v3, v3s
)

var (
	mI = map[int]string{
		v1: v1s,
		v2: v2s,
		v3: v3s,
		v4: v4s,
	}
	mS = func() map[string]int {
		result := make(map[string]int, len(mI))
		for k, v := range mI {
			result[v] = k
		}
		return result
	}()
)

func BenchmarkMap(b *testing.B) {
	b.Run("int to string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = mI[lookupI]
		}
	})
	b.Run("string to int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = mS[lookupS]
		}
	})
}

func BenchmarkSwitch(b *testing.B) {
	b.Run("int to string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			switch lookupI {
			case v1:
			case v2:
			case v3:
			case v4:
			}
		}
	})
	b.Run("string to int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			switch lookupS {
			case v1s:
			case v2s:
			case v3s:
			case v4s:
			}
		}
	})
}

func BenchmarkIf(b *testing.B) {
	lookupI, lookupS := v3, v3s
	b.Run("int to string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if lookupI == v1 {

			} else if lookupI == v2 {

			} else if lookupI == v3 {

			} else if lookupI == v4 {

			}
		}
	})
	b.Run("string to int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if lookupS == v1s {

			} else if lookupS == v2s {

			} else if lookupS == v3s {

			} else if lookupS == v4s {

			}
		}
	})
}

func BenchmarkString(b *testing.B) {
	searchIndex := 3
	b.Run("enumer method", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			StringPill(searchIndex)
		}
	})
	b.Run("switch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			StringSwitch(searchIndex)
		}
	})
	b.Run("switch2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			StringSwitch2(16)
		}
	})
}

const _PillName = "PlaceboAspirinIbuprofenParacetamol"

var _PillIndex = [...]uint8{0, 7, 14, 23, 34}

func StringPill(enumIndex int) string {
	if enumIndex < 0 || enumIndex >= len(_PillIndex)-1 {
		return fmt.Sprintf("Pill(%d)", enumIndex)
	}
	return _PillName[_PillIndex[enumIndex]:_PillIndex[enumIndex+1]]
}

func StringSwitch(enumIndex int) string {
	switch enumIndex {
	case 0:
		return "Placebo"
	case 1:
		return "Aspirin"
	case 2:
		return "Ibuprofen"
	case 3:
		return "Paracetamol"
	}
	return fmt.Sprintf("Pill(%d)", enumIndex)
}

func StringSwitch2(enumIndex int) string {
	switch enumIndex {
	case 0:
		return "Enum3Value0"
	case 1:
		return "Enum3Value1"
	case 2:
		return "Enum3Value2"
	case 3:
		return "Enum3Value3"
	case 4:
		return "Enum3Value4"
	case 5:
		return "Enum3Value5"
	case 6:
		return "Enum3Value6"
	case 7:
		return "Enum3Value7"
	case 8:
		return "Enum3Value8"
	case 9:
		return "Enum3Value9"
	case 10:
		return "Enum3Value10"
	case 11:
		return "Enum3Value11"
	case 12:
		return "Enum3Value12"
	case 13:
		return "Enum3Value13"
	case 14:
		return "Enum3Value14"
	case 15:
		return "Enum3Value15"
	case 16:
		return "Enum3Value16"
	default:
		return fmt.Sprintf("UndefinedMyEnum3:%d", enumIndex)
	}
}
