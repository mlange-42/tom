package api_test

import (
	"testing"

	"github.com/mlange-42/tom/api"
	"github.com/stretchr/testify/assert"
)

func TestDirection(t *testing.T) {
	assert.Equal(t, "N", api.Direction(0))
	assert.Equal(t, "E", api.Direction(90))
	assert.Equal(t, "S", api.Direction(180))
	assert.Equal(t, "W", api.Direction(270))
	assert.Equal(t, "N", api.Direction(360))

	assert.Equal(t, "N", api.Direction(22))
	assert.Equal(t, "NE", api.Direction(23))
}
