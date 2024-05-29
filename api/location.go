package api

import (
	"context"
	"fmt"
)

type Location struct {
	Lat      float64
	Lon      float64
	TimeZone string
}

func GetLocation(loc string, locations map[string]Location) (Location, error) {
	coords, ok := locations[loc]
	if !ok {
		client := NewClient(Geocoding)
		opt := GeoOptions{
			Name: loc,
		}
		result, err := client.Get(context.Background(), &opt)
		if err != nil {
			return Location{}, err
		}
		parsed, err := ParseGeo(result)
		if err != nil {
			return Location{}, err
		}
		if len(parsed.Results) == 0 {
			return Location{}, fmt.Errorf("location not found: '%s'", loc)
		}
		coords = Location{
			Lat:      parsed.Results[0].Latitude,
			Lon:      parsed.Results[0].Longitude,
			TimeZone: parsed.Results[0].TimeZone,
		}
	}
	return coords, nil
}
