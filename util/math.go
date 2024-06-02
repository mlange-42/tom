package util

import "cmp"

func Clamp[E cmp.Ordered](x, mn, mx E) E {
	return max(min(x, mx), mn)
}
