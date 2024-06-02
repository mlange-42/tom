package render_test

import (
	"testing"

	"github.com/mlange-42/tom/render"
	"github.com/stretchr/testify/assert"
)

func TestCanvas(t *testing.T) {
	c := render.NewCanvas(10, 4)

	sx, sy := c.PixelSize()

	assert.Equal(t, 20, sx)
	assert.Equal(t, 16, sy)

	c.Set(2, 0, true)
	c.Set(3, 0, true)
	c.Set(2, 1, true)

	runes := c.Runes()
	text := c.String()

	assert.Equal(t, '⠋', runes[0][1])
	assert.Equal(t, '⠋', []rune(text)[1])
}
