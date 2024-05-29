package api

import (
	"math"
)

var directions = []string{
	"N",
	"NNE",
	"NE",
	"ENE",
	"E",
	"ESE",
	"SE",
	"SSE",
	"S",
	"SSW",
	"SW",
	"WSW",
	"W",
	"WNW",
	"NW",
	"WWN",
}

func Direction(degrees float64) string {
	offset := 11.25
	step := 22.5
	idx := int(math.Mod(degrees+offset, 360) / step)
	return directions[idx]
}
