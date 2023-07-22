package math

// Integer is a type constraint that encompasses all integer types.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Float is a type constraint that encompasses all float types.
type Float interface {
	~float32 | ~float64
}

// Number is a type constraint that encompasses all number types.
type Number interface {
	Integer | Float
}

// Max returns the max value of any number.
func Max[T Number](a T, b T, vals ...T) T {
	if len(vals) == 0 {
		if b > a {
			return b
		}
		return a
	}

	if b > a {
		a = b
	}
	for _, v := range vals {
		if v > a {
			a = v
		}
	}
	return a
}

// Min returns the min value of any number.
func Min[T Number](a T, b T, vals ...T) T {
	if len(vals) == 0 {
		if b < a {
			return b
		}
		return a
	}

	if b < a {
		a = b
	}
	for _, v := range vals {
		if v < a {
			a = v
		}
	}
	return a
}

// Constrain ensures a number is between a given min and max.
func Constrain[T Number](val T, min T, max T) T {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
