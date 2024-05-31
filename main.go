package main

import (
	"log"
	"os"
	"strings"

	"github.com/mlange-42/tom/api"
	"github.com/mlange-42/tom/app"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("please specify a location")
	}

	loc := strings.ToLower(os.Args[1])
	coords, err := api.GetLocation(loc)
	if err != nil {
		log.Fatal(err)
	}

	meteo, err := api.GetMeteo(coords)
	if err != nil {
		log.Fatal(err)
	}

	a := app.New(strings.ToTitle(loc), meteo)
	a.Run()
}
