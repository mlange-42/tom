package util

// Repeat creates a slice og the given length, filled with the given value.
func Repeat[T any](value T, length int) []T {
	sl := make([]T, length)
	for i := range sl {
		sl[i] = value
	}
	return sl
}
