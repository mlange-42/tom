package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
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
	GenerationTime_ms float64
	Current           CurrentWeather

	Hourly     map[string][]float64
	HourlyTime []time.Time

	Daily     map[string][]float64
	DailyTime []time.Time
}

type meteoResultJs struct {
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
	current.Time, err = time.Parse(timeLayout, m.Current["time"].(string))
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

	return &MeteoResult{
		GenerationTime_ms: m.GenerationTime_ms,
		Current:           current,
	}, nil
}
