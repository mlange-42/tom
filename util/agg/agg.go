package agg

import (
	"math"
	"time"
)

func AggregateTime(t []time.Time, step int, point int) []time.Time {
	ln := len(t) / step

	times := make([]time.Time, ln)

	for i := 0; i < ln; i++ {
		idx := i*step + point
		times[i] = t[idx]
	}

	return times
}

type Aggregator interface {
	Aggregate(data []float64, step int, point int) []float64
}

type Point struct{}

func (a *Point) Aggregate(data []float64, step int, point int) []float64 {
	ln := len(data) / step

	values := make([]float64, ln)

	for i := 0; i < ln; i++ {
		idx := i*step + point
		values[i] = data[idx]
	}

	return values
}

type Max struct{}

func (a *Max) Aggregate(data []float64, step int, point int) []float64 {
	ln := len(data) / step

	values := make([]float64, ln)

	for i := 0; i < ln; i++ {
		idx := i * step
		start := idx
		end := idx + step
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

		values[i] = mx
	}

	return values
}

type Sum struct{}

func (a *Sum) Aggregate(data []float64, step int, point int) []float64 {
	ln := len(data) / step

	values := make([]float64, ln)

	for i := 0; i < ln; i++ {
		idx := i * step
		start := idx
		end := idx + step
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

		values[i] = sum
	}

	return values
}

type Avg struct{}

func (a *Avg) Aggregate(data []float64, step int, point int) []float64 {
	ln := len(data) / step

	values := make([]float64, ln)

	for i := 0; i < ln; i++ {
		idx := i * step
		start := idx
		end := idx + step
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

		values[i] = sum / float64(cnt)
	}

	return values
}

type ModeOrPoint struct{}

func (a *ModeOrPoint) Aggregate(data []float64, step int, point int) []float64 {
	ln := len(data) / step

	values := make([]float64, ln)
	tempValues := map[float64]int{}

	for i := 0; i < ln; i++ {
		idx := i * step
		start := idx
		end := idx + step
		if start < 0 {
			start = 0
		}
		if end > len(data) {
			end = len(data)
		}

		vByValue := data[idx+point]
		for j := start; j < end; j++ {
			v := data[j]
			if cnt, ok := tempValues[v]; ok {
				tempValues[v] = cnt + 1
			} else {
				tempValues[v] = 1
			}
			if v > vByValue {
				vByValue = v
			}
		}
		vByCnt := data[idx+point]
		cntMax := 1
		for v, cnt := range tempValues {
			if cnt > cntMax {
				vByCnt = v
			}
		}
		if cntMax <= 1 {
			values[i] = vByValue
		} else {
			values[i] = vByCnt
		}

		clear(tempValues)
	}

	return values
}
