package api

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mlange-42/tom/config"
)

const DefaultUserAgent = "TOM-terminal-open-meteo"

type API string

const (
	OpenMeteo API = "https://api.open-meteo.com/v1"
	Geocoding API = "https://geocoding-api.open-meteo.com/v1"
)

type Client interface {
	Get(ctx context.Context, req config.Options) ([]byte, error)
}

type openMeteoClient struct {
	URL       string
	UserAgent string
	Client    *http.Client
}

func NewClient(api API) Client {
	return &openMeteoClient{
		URL:       string(api),
		UserAgent: DefaultUserAgent,
		Client:    http.DefaultClient,
	}
}

func (c *openMeteoClient) Get(ctx context.Context, opt config.Options) ([]byte, error) {
	url := opt.ToURL(c.URL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%s - %s", res.Status, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetLocation(location string) (config.Location, error) {
	cached, err := config.LoadLocations()
	if err != nil {
		return config.Location{}, err
	}

	coords, ok := cached[location]
	if !ok {
		locations, err := GetLocations(location)
		if err != nil {
			return config.Location{}, err
		}

		coords = config.Location{
			Lat:      locations[0].Latitude,
			Lon:      locations[0].Longitude,
			TimeZone: locations[0].TimeZone,
		}

		cached[location] = coords
		err = config.SaveLocations(cached)
		if err != nil {
			return config.Location{}, err
		}
	}

	return coords, nil
}

func GetLocations(loc string) ([]config.GeoResultEntry, error) {
	client := NewClient(Geocoding)
	opt := config.GeoOptions{
		Name: loc,
	}
	result, err := client.Get(context.Background(), &opt)
	if err != nil {
		return nil, err
	}
	parsed, err := config.ParseGeo(result)
	if err != nil {
		return nil, err
	}

	locations := []config.GeoResultEntry{}
	for _, e := range parsed.Results {
		if e.TimeZone != "" {
			locations = append(locations, e)
		}
	}
	return locations, nil
}

func GetMeteo(loc config.Location) (*config.MeteoResult, error) {
	client := NewClient(OpenMeteo)

	opt := config.ForecastOptions{
		Location:       loc,
		Days:           7,
		CurrentMetrics: config.DefaultCurrentMetrics,
		HourlyMetrics:  config.DefaultHourlyMetrics,
		DailyMetrics:   config.DefaultDailyMetrics,
	}

	result, err := client.Get(context.Background(), &opt)
	if err != nil {
		return nil, err
	}

	return config.ParseMeteo(result, &opt)
}
