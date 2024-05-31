package api

import (
	"context"
	"fmt"

	"github.com/mlange-42/tom/config"
)

func getLocation(loc string, locations map[string]config.Location) (config.Location, error) {
	coords, ok := locations[loc]
	if !ok {
		client := NewClient(Geocoding)
		opt := config.GeoOptions{
			Name: loc,
		}
		result, err := client.Get(context.Background(), &opt)
		if err != nil {
			return config.Location{}, err
		}
		parsed, err := config.ParseGeo(result)
		if err != nil {
			return config.Location{}, err
		}
		if len(parsed.Results) == 0 {
			return config.Location{}, fmt.Errorf("location not found: '%s'", loc)
		}
		coords = config.Location{
			Lat:      parsed.Results[0].Latitude,
			Lon:      parsed.Results[0].Longitude,
			TimeZone: parsed.Results[0].TimeZone,
		}
	}
	return coords, nil
}
