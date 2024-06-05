package render

import (
	"math"

	"github.com/mlange-42/tom/util"
)

type Chart struct {
	canvas *Canvas
}

func NewChart(w, h int) *Chart {
	return &Chart{
		canvas: NewCanvas(w, h),
	}
}

func (c *Chart) Series(data []float64, bars bool) (min float64, max float64) {
	min, max = util.MinMax(data)
	if min > 0 {
		min = 0
	}
	if max == min {
		max = min + 1
	}
	_, height := c.canvas.PixelSize()
	height -= 1

	yZero := int(math.Round(float64(height) * (-min) / (max - min)))
	yZero = util.Clamp(yZero, 0, height)

	for x, v := range data {
		y := int(math.Round(float64(height) * (v - min) / (max - min)))
		c.canvas.Set(x, height-y, true)
		if bars {
			if v >= 0 {
				for yy := yZero; yy < y; yy++ {
					c.canvas.Set(x, height-yy, true)
				}
			} else {
				for yy := yZero; yy > y; yy-- {
					c.canvas.Set(x, height-yy, true)
				}
			}
		} else {
			//if yZero != 0 {
			c.canvas.Set(x, height-yZero, true)
			//}
		}
	}

	return
}

func (c *Chart) VLine(x int, dash int, invert bool) {
	_, height := c.canvas.PixelSize()
	if dash < 1 {
		dash = 1
	}
	if invert {
		for y := 0; y < height; y += dash {
			c.canvas.Set(x, y, !c.canvas.Get(x, y))
		}
	} else {
		for y := 0; y < height; y += dash {
			c.canvas.Set(x, y, true)
		}
	}
}

func (c *Chart) Runes() [][]rune {
	return c.canvas.Runes()
}

func (c *Chart) String() string {
	return c.canvas.String()
}
