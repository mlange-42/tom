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

		var selected *config.GeoResultEntry
		for i := range parsed.Results {
			e := &parsed.Results[i]
			if e.TimeZone != "" {
				selected = e
				break
			}
		}

		if selected == nil {
			return config.Location{}, fmt.Errorf("location not found: '%s'", loc)
		}
		coords = config.Location{
			Lat:      selected.Latitude,
			Lon:      selected.Longitude,
			TimeZone: selected.TimeZone,
		}
	}
	return coords, nil
}
