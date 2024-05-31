package config_test

import (
	"testing"

	"github.com/mlange-42/tom/config"
	"github.com/stretchr/testify/assert"
)

func TestDirection(t *testing.T) {
	assert.Equal(t, "N", config.Direction(0))
	assert.Equal(t, "E", config.Direction(90))
	assert.Equal(t, "S", config.Direction(180))
	assert.Equal(t, "W", config.Direction(270))
	assert.Equal(t, "N", config.Direction(360))

	assert.Equal(t, "N", config.Direction(22))
	assert.Equal(t, "NE", config.Direction(23))
}
