package main

import (
	"log"
	"os"
	"strings"

	"github.com/mlange-42/tom/app"
	"github.com/mlange-42/tom/config"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("please specify a location")
	}
	location := strings.ToLower(os.Args[1])

	cached, err := config.LoadLocations()
	if err != nil {
		log.Fatal(err)
	}

	coords, ok := cached[location]
	if ok {
		a := app.New(strings.ToTitle(location), coords)
		a.Run()
		return
	}
	a := app.NewLocationDialog(location, cached)
	a.Run()
}
