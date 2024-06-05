package render

import (
	"strings"
)

const brailleWidth = 2
const brailleHeight = 4

type Braille [2][4]uint8

func (b *Braille) Rune() rune {
	littleEndian := [8]uint8{
		b[0][0], b[0][1], b[0][2],
		b[1][0], b[1][1], b[1][2],
		b[0][3], b[1][3],
	}
	var v = 0
	for i, x := range littleEndian {
		v += int(x) << i
	}
	return rune(v) + '\u2800'
}

type Canvas struct {
	width  int
	height int
	grid   [][]Braille
}

func NewCanvas(w, h int) *Canvas {
	grid := make([][]Braille, h)
	for i := range grid {
		grid[i] = make([]Braille, w)
	}
	return &Canvas{
		width:  w,
		height: h,
		grid:   grid,
	}
}

func (c *Canvas) Size() (int, int) {
	return c.width, c.height
}

func (c *Canvas) PixelSize() (int, int) {
	return c.width * brailleWidth, c.height * brailleHeight
}

func (c *Canvas) Set(x, y int, v bool) {
	col, xx := x/brailleWidth, x%brailleWidth
	row, yy := y/brailleHeight, y%brailleHeight

	b := &c.grid[row][col]
	if v {
		b[xx][yy] = 1
	} else {
		b[xx][yy] = 0
	}
}

func (c *Canvas) Get(x, y int) bool {
	col, xx := x/brailleWidth, x%brailleWidth
	row, yy := y/brailleHeight, y%brailleHeight

	b := &c.grid[row][col]
	return b[xx][yy] > 0
}

func (c *Canvas) Runes() [][]rune {
	runes := make([][]rune, c.height)

	for i, row := range c.grid {
		rowRunes := make([]rune, len(row))
		for j, r := range row {
			rowRunes[j] = r.Rune()
		}
		runes[i] = rowRunes
	}
	return runes
}

func (c *Canvas) String() string {
	lines := make([]string, c.height)

	for i, row := range c.grid {
		builder := strings.Builder{}
		for _, r := range row {
			builder.WriteRune(r.Rune())
		}
		lines[i] = builder.String()
	}
	return strings.Join(lines, "\n")
}
