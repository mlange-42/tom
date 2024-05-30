package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mlange-42/tom/util/agg"
)

type GeoResult struct {
	Results []GeoResultEntry
}

type GeoResultEntry struct {
	Name      string
	Latitude  float64
	Longitude float64
	Country   string
	Admin1    string
	Elevation float64
	TimeZone  string
}

type MeteoResult struct {
	Location Location

	GenerationTime_ms float64
	Current           CurrentWeather

	Hourly     map[string][]float64
	HourlyTime []time.Time

	Daily     map[string][]float64
	DailyTime []time.Time

	ThreeHourly     map[string][]float64
	ThreeHourlyTime []time.Time

	SixHourly     map[string][]float64
	SixHourlyTime []time.Time
}

func (r *MeteoResult) GetCurrent(key CurrentMetric) float64 {
	return r.Current.Values[string(key)]
}

func (r *MeteoResult) GetDaily(key CurrentMetric) []float64 {
	return r.Daily[string(key)]
}

func (r *MeteoResult) GetHourly(key HourlyMetric) []float64 {
	return r.Hourly[string(key)]
}

func (r *MeteoResult) GetThreeHourly(key HourlyMetric) []float64 {
	return r.ThreeHourly[string(key)]
}

func (r *MeteoResult) GetSixHourly(key HourlyMetric) []float64 {
	return r.SixHourly[string(key)]
}

type meteoResultJs struct {
	Latitude          float64
	Longitude         float64
	TimeZone          string
	GenerationTime_ms float64
	Current           map[string]any             `json:"current"`
	Hourly            map[string]json.RawMessage `json:"hourly"`
	Daily             map[string]json.RawMessage `json:"daily"`
}

type CurrentWeather struct {
	Time   time.Time
	Values map[string]float64
}

func ParseGeo(data []byte) (*GeoResult, error) {
	r := GeoResult{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

func ParseMeteo(data []byte, opt *ForecastOptions) (*MeteoResult, error) {
	m := meteoResultJs{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&m); err != nil {
		return nil, err
	}

	var err error
	current := CurrentWeather{Values: map[string]float64{}}
	current.Time, err = time.Parse(DateTimeLayout, m.Current["time"].(string))
	if err != nil {
		return nil, err
	}
	for _, key := range opt.CurrentMetrics {
		v, ok := m.Current[string(key)]
		if !ok {
			return nil, fmt.Errorf("metric '%s' not in results", string(key))
		}
		f, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("can't convert '%s' to float64", v)
		}
		current.Values[string(key)] = f
	}

	hourlyTime, hourly, err := parseHourly(&m, opt.HourlyMetrics)
	if err != nil {
		return nil, err
	}

	dailyTime, daily, err := parseDaily(&m, opt.DailyMetrics)
	if err != nil {
		return nil, err
	}

	threeHourlyTime := agg.AggregateTime(hourlyTime, 3)
	threeHourly := map[string][]float64{}
	for _, key := range opt.HourlyMetrics {
		threeHourly[string(key)] = aggregators[key].Aggregate(hourly[string(key)], 3, 1, 1)
	}

	sixHourlyTime := agg.AggregateTime(hourlyTime, 6)
	sixHourly := map[string][]float64{}
	for _, key := range opt.HourlyMetrics {
		sixHourly[string(key)] = aggregators[key].Aggregate(hourly[string(key)], 6, 1, 4)
	}

	return &MeteoResult{
		Location: Location{
			Lat:      m.Latitude,
			Lon:      m.Longitude,
			TimeZone: m.TimeZone,
		},
		GenerationTime_ms: m.GenerationTime_ms,
		Current:           current,
		Hourly:            hourly,
		HourlyTime:        hourlyTime,
		Daily:             daily,
		DailyTime:         dailyTime,
		ThreeHourly:       threeHourly,
		ThreeHourlyTime:   threeHourlyTime,
		SixHourly:         sixHourly,
		SixHourlyTime:     sixHourlyTime,
	}, nil
}

func parseHourly(m *meteoResultJs, metrics []HourlyMetric) ([]time.Time, map[string][]float64, error) {
	hourly := map[string][]float64{}
	rawTime, ok := m.Hourly["time"]
	if !ok {
		return nil, nil, fmt.Errorf("no time not in results")
	}
	timeStr := []string{}
	err := json.Unmarshal(rawTime, &timeStr)
	if err != nil {
		return nil, nil, err
	}
	hourlyTime := make([]time.Time, len(timeStr))
	for i, v := range timeStr {
		hourlyTime[i], err = time.Parse(DateTimeLayout, v)
		if err != nil {
			return nil, nil, err
		}
	}

	for _, key := range metrics {
		v, ok := m.Hourly[string(key)]
		if !ok {
			return nil, nil, fmt.Errorf("metric '%s' not in results", string(key))
		}
		data := []float64{}
		err := json.Unmarshal(v, &data)
		if err != nil {
			return nil, nil, err
		}
		hourly[string(key)] = data
	}
	return hourlyTime, hourly, nil
}

func parseDaily(m *meteoResultJs, metrics []DailyMetric) ([]time.Time, map[string][]float64, error) {
	daily := map[string][]float64{}
	rawTime, ok := m.Daily["time"]
	if !ok {
		return nil, nil, fmt.Errorf("no time not in results")
	}
	timeStr := []string{}
	err := json.Unmarshal(rawTime, &timeStr)
	if err != nil {
		return nil, nil, err
	}
	dailyTime := make([]time.Time, len(timeStr))
	for i, v := range timeStr {
		dailyTime[i], err = time.Parse(DateLayout, v)
		if err != nil {
			return nil, nil, err
		}
	}

	for _, key := range metrics {
		v, ok := m.Daily[string(key)]
		if !ok {
			return nil, nil, fmt.Errorf("metric '%s' not in results", string(key))
		}
		data := []float64{}
		err := json.Unmarshal(v, &data)
		if err != nil {
			return nil, nil, err
		}
		daily[string(key)] = data
	}
	return dailyTime, daily, nil
}
