package agg

import (
	"math"
	"time"
)

type Aggregator interface {
	Aggregate(data []float64, step int, past int, future int) ([]time.Time, []float64)
}

type Point struct{}

func (a *Point) Aggregate(t []time.Time, data []float64, step int, past int, future int) ([]time.Time, []float64) {
	ln := len(data) / step

	times := make([]time.Time, ln)
	values := make([]float64, ln)

	for i := 0; i < ln; i++ {
		idx := i * step
		times[i] = t[idx]
		values[i] = data[idx]
	}

	return times, values
}

type Max struct{}

func (a *Max) Aggregate(t []time.Time, data []float64, step int, past int, future int) ([]time.Time, []float64) {
	ln := len(data) / step

	times := make([]time.Time, ln)
	values := make([]float64, ln)

	for i := 0; i < ln; i++ {
		idx := i * step
		start := idx - past
		end := idx + future + 1
		if start < 0 {
			start = 0
		}
		if end > len(data) {
			end = len(data)
		}

		mx := math.Inf(-1)
		for j := start; j < end; j++ {
			if data[j] > mx {
				mx = data[j]
			}
		}

		times[i] = t[idx]
		values[i] = mx
	}

	return times, values
}

type Sum struct{}

func (a *Sum) Aggregate(t []time.Time, data []float64, step int, past int, future int) ([]time.Time, []float64) {
	ln := len(data) / step

	times := make([]time.Time, ln)
	values := make([]float64, ln)

	for i := 0; i < ln; i++ {
		idx := i * step
		start := idx - past
		end := idx + future + 1
		if start < 0 {
			start = 0
		}
		if end > len(data) {
			end = len(data)
		}

		sum := 0.0
		for j := start; j < end; j++ {
			sum += data[j]
		}

		times[i] = t[idx]
		values[i] = sum
	}

	return times, values
}

type Avg struct{}

func (a *Avg) Aggregate(t []time.Time, data []float64, step int, past int, future int) ([]time.Time, []float64) {
	ln := len(data) / step

	times := make([]time.Time, ln)
	values := make([]float64, ln)

	for i := 0; i < ln; i++ {
		idx := i * step
		start := idx - past
		end := idx + future + 1
		if start < 0 {
			start = 0
		}
		if end > len(data) {
			end = len(data)
		}

		sum := 0.0
		cnt := 0
		for j := start; j < end; j++ {
			sum += data[j]
			cnt++
		}

		times[i] = t[idx]
		values[i] = sum / float64(cnt)
	}

	return times, values
}
