package api

import (
	"math"
)

var directions = []string{
	"N",
	"NE",
	"E",
	"SE",
	"S",
	"SW",
	"W",
	"NW",
}

func Direction(degrees float64) string {
	step := 360.0 / float64(len(directions)) // 45.0°
	offset := step / 2.0
	idx := int(math.Mod(degrees+offset, 360) / step)
	return directions[idx]
}
