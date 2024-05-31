package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/mlange-42/tom/api"
	"github.com/mlange-42/tom/config"
	"github.com/stretchr/testify/assert"
)

func TestGeoClient(t *testing.T) {
	client := api.NewClient(api.Geocoding)

	opt := config.GeoOptions{
		Name: "Berlin",
	}

	result, err := client.Get(context.Background(), &opt)
	assert.Nil(t, err)

	parsed, err := config.ParseGeo(result)
	assert.Nil(t, err)

	assert.Equal(t, 10, len(parsed.Results))

	assert.Equal(t, "Germany", parsed.Results[0].CountryCode)
	assert.Equal(t, "Land Berlin", parsed.Results[0].Admin1)
	assert.Equal(t, "Europe/Berlin", parsed.Results[0].TimeZone)
	assert.Greater(t, parsed.Results[0].Latitude, 0.0)
	assert.Greater(t, parsed.Results[0].Longitude, 0.0)
	assert.Greater(t, parsed.Results[0].Elevation, 0.0)
}

func TestMeteoClient(t *testing.T) {
	client := api.NewClient(api.OpenMeteo)

	opt := config.ForecastOptions{
		Location: config.Location{
			Lat:      52.5,
			Lon:      13.4,
			TimeZone: "Europe/Berlin",
		},
		Days:           3,
		CurrentMetrics: []config.CurrentMetric{config.CurrentWindSpeed, config.CurrentRH},
		HourlyMetrics:  []config.HourlyMetric{config.HourlyTemp},
		DailyMetrics:   []config.DailyMetric{config.DailyMaxTemp, config.DailyMinTemp},
	}

	result, err := client.Get(context.Background(), &opt)
	assert.Nil(t, err)

	parsed, err := config.ParseMeteo(result, &opt)
	assert.Nil(t, err)

	assert.Greater(t, parsed.GenerationTime_ms, 0.0)
	assert.Greater(t, parsed.Current.Time, time.Time{})

	assert.Greater(t, parsed.Current.Values[string(config.CurrentWindSpeed)], 0.0)
	assert.Greater(t, parsed.Current.Values[string(config.CurrentRH)], 0.0)

	assert.Greater(t, len(parsed.HourlyTime), 0)
	temp, ok := parsed.Hourly[string(config.HourlyTemp)]
	assert.True(t, ok)
	assert.Greater(t, len(temp), 0)

	assert.Greater(t, len(parsed.DailyTime), 0)
	maxTemp, ok := parsed.Daily[string(config.DailyMaxTemp)]
	assert.True(t, ok)
	assert.Greater(t, len(maxTemp), 0)

	assert.Greater(t, parsed.Location.Lat, 0.0)
	assert.Greater(t, parsed.Location.Lon, 0.0)
	assert.Equal(t, "Europe/Berlin", parsed.Location.TimeZone)

	assert.Equal(t, len(parsed.HourlyTime)/3, len(parsed.ThreeHourlyTime))
	assert.Equal(t, len(parsed.HourlyTime)/6, len(parsed.SixHourlyTime))

	assert.Equal(t, len(parsed.Hourly[string(config.HourlyTemp)])/3, len(parsed.ThreeHourly[string(config.HourlyTemp)]))
	assert.Equal(t, len(parsed.Hourly[string(config.HourlyTemp)])/6, len(parsed.SixHourly[string(config.HourlyTemp)]))
}
