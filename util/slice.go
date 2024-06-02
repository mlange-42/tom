package util

import "cmp"

// Repeat creates a slice og the given length, filled with the given value.
func Repeat[T any](value T, length int) []T {
	sl := make([]T, length)
	for i := range sl {
		sl[i] = value
	}
	return sl
}

func MinMax[S ~[]E, E cmp.Ordered](x S) (E, E) {
	if len(x) < 1 {
		panic("slices.Max: empty list")
	}
	mn := x[0]
	mx := x[0]
	for i := 1; i < len(x); i++ {
		mn = min(mn, x[i])
		mx = max(mx, x[i])
	}
	return mn, mx
}
