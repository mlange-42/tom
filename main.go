package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mlange-42/tom/api"
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

	fmt.Println(strings.ToTitle(loc), coords)
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
