package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/mlange-42/tom/api"
	"github.com/stretchr/testify/assert"
)

func TestGeoClient(t *testing.T) {
	client := api.NewClient(api.Geocoding)

	opt := api.GeoOptions{
		Name: "Berlin",
	}

	result, err := client.Get(context.Background(), &opt)
	assert.Nil(t, err)

	parsed, err := api.ParseGeo(result)
	assert.Nil(t, err)

	assert.Equal(t, 10, len(parsed.Results))

	assert.Equal(t, "Germany", parsed.Results[0].Country)
	assert.Equal(t, "Land Berlin", parsed.Results[0].Admin1)
	assert.Equal(t, "Europe/Berlin", parsed.Results[0].TimeZone)
	assert.Greater(t, parsed.Results[0].Latitude, 0.0)
	assert.Greater(t, parsed.Results[0].Longitude, 0.0)
	assert.Greater(t, parsed.Results[0].Elevation, 0.0)
}

func TestMeteoClient(t *testing.T) {
	client := api.NewClient(api.OpenMeteo)

	opt := api.ForecastOptions{
		Location: api.Location{
			Lat: 13.4,
			Lon: 52.5,
		},
		Timezone:       "Europe/Berlin",
		Days:           2,
		CurrentMetrics: []api.CurrentMetric{api.CurrentWindSpeed, api.CurrentRH},
		HourlyMetrics:  []api.HourlyMetric{api.HourlyTemp},
		DailyMetrics:   []api.DailyMetric{api.DailyMaxTemp},
	}

	result, err := client.Get(context.Background(), &opt)
	assert.Nil(t, err)

	parsed, err := api.ParseMeteo(result, &opt)
	assert.Nil(t, err)

	assert.Greater(t, parsed.GenerationTime_ms, 0.0)
	assert.Greater(t, parsed.Current.Time, time.Time{})

	assert.Greater(t, parsed.Current.Values[string(api.CurrentWindSpeed)], 0.0)
	assert.Greater(t, parsed.Current.Values[string(api.CurrentRH)], 0.0)
}
