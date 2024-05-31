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
	return directions[directionIndex(degrees)]
}

func directionIndex(degrees float64) int {
	step := 360.0 / float64(len(directions)) // 45.0°
	offset := step / 2.0
	idx := int(math.Mod(degrees+offset, 360) / step)
	return idx
}
