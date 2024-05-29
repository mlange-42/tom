package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/mlange-42/tom/api"
	"github.com/mlange-42/tom/app"
	"github.com/mlange-42/tom/config"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("please specify a location")
	}

	loc := strings.ToLower(os.Args[1])
	coords, err := GetLocation(loc)
	if err != nil {
		log.Fatal(err)
	}

	meteo, err := GetMeteo(coords)
	if err != nil {
		log.Fatal(err)
	}

	a := app.New(strings.ToTitle(loc), meteo)
	a.Run(app.TCellTerminal)
}

func GetLocation(location string) (api.Location, error) {
	locations, err := config.LoadLocations()
	if err != nil {
		return api.Location{}, err
	}
	coords, err := api.GetLocation(location, locations)
	if err != nil {
		return api.Location{}, err
	}
	locations[location] = coords
	err = config.SaveLocations(locations)
	if err != nil {
		return api.Location{}, err
	}
	return coords, nil
}

func GetMeteo(loc api.Location) (*api.MeteoResult, error) {
	client := api.NewClient(api.OpenMeteo)

	opt := api.ForecastOptions{
		Location: loc,
		Days:     7,
		CurrentMetrics: []api.CurrentMetric{
			api.CurrentWeatherCode,
			api.CurrentTemp,
			api.CurrentPrecip,
			api.CurrentCloudCover,
			api.CurrentWindSpeed,
			api.CurrentWindDir,
			api.CurrentRH,
		},
		HourlyMetrics: []api.HourlyMetric{
			api.HourlyWeatherCode,
			api.HourlyTemp,
			api.HourlyPrecip,
			api.HourlyPrecipProb,
			api.HourlyRH,
			api.HourlyCloudCover,
			api.HourlyWindSpeed,
			api.HourlyWindDir,
		},
		DailyMetrics: []api.DailyMetric{
			api.DailyWeatherCode,
			api.DailyMinTemp,
			api.DailyMaxTemp,
			api.DailyPrecip,
			api.DailyPrecipProb,
			api.DailySunshine,
			api.DailyWindSpeed,
			api.DailyWindDir,
		},
	}

	result, err := client.Get(context.Background(), &opt)
	if err != nil {
		return nil, err
	}

	return api.ParseMeteo(result, &opt)
}
