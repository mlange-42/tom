package agg_test

import (
	"testing"
	"time"

	"github.com/mlange-42/tom/util/agg"
	"github.com/stretchr/testify/assert"
)

func TestAggTime(t *testing.T) {
	times, _ := createData()

	tt := agg.AggregateTime(times, 3, 0)
	assert.Equal(t, 8, len(tt))

	assert.Equal(t, time.Date(2000, 1, 1, 3, 0, 0, 0, time.UTC), tt[1])
}

func TestAggPoint(t *testing.T) {
	_, data := createData()

	a := agg.Point{}

	v := a.Aggregate(data, 3, 1)
	assert.Equal(t, 8, len(v))

	assert.Equal(t, []float64{1, 4, 7, 10, 13, 16, 19, 22}, v)
}

func TestAggMax(t *testing.T) {
	_, data := createData()

	a := agg.Max{}

	v := a.Aggregate(data, 3, 1)
	assert.Equal(t, 8, len(v))

	assert.Equal(t, []float64{2, 5, 8, 11, 14, 17, 20, 23}, v)
}

func TestAggSum(t *testing.T) {
	_, data := createData()

	a := agg.Sum{}

	v := a.Aggregate(data, 3, 1)
	assert.Equal(t, 8, len(v))

	assert.Equal(t, []float64{0 + 1 + 2, 3 + 4 + 5, 6 + 7 + 8, 9 + 10 + 11, 12 + 13 + 14, 15 + 16 + 17, 18 + 19 + 20, 21 + 22 + 23}, v)
}

func TestAggAvg(t *testing.T) {
	_, data := createData()

	a := agg.Avg{}

	v := a.Aggregate(data, 3, 1)
	assert.Equal(t, 8, len(v))

	assert.Equal(t, []float64{1, 4, 7, 10, 13, 16, 19, 22}, v)
}

func createData() ([]time.Time, []float64) {
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	times := make([]time.Time, 24)
	data := make([]float64, 24)

	for i := range times {
		times[i] = start.Add(time.Duration(i) * time.Hour)
		data[i] = float64(i)
	}
	return times, data
}
