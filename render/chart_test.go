package render_test

import (
	"testing"

	"github.com/mlange-42/tom/render"
	"github.com/stretchr/testify/assert"
)

func TestChart(t *testing.T) {
	c := render.NewChart(5, 4)

	data := []float64{0, 1, 2, 6, 12, 8, 4, 1, -2, -4}

	c.Series(data, true)

	assert.Equal(t,
		`⠀⠀⡇⠀⠀
⠀⢠⣿⡀⠀
⣠⣾⣿⣧⣀
⠀⠀⠀⠀⢻`, c.String())
}
