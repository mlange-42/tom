package data

import (
	_ "embed"
	"log"
	"strings"

	"github.com/mlange-42/tom/config"
	"gopkg.in/yaml.v3"
)

//go:embed weather_codes.yml
var weatherCodeBytes []byte

//go:embed day_layout.txt
var dayLayout string

var DayLayout [][]rune

var WeatherCodes map[int]WeatherCode

type WeatherCode struct {
	Name   string
	Symbol [][]rune
	Colors [][]uint8
}

type weatherCode struct {
	Name   string
	Symbol string
	Color  string
}

var ColorIDs map[rune]uint8

func init() {
	ColorIDs = map[rune]uint8{}
	for i, c := range config.Colors {
		ColorIDs[c.Rune] = uint8(i)
	}

	l := strings.Replace(dayLayout, "\r\n", "\n", -1)
	lines := strings.Split(l, "\n")

	DayLayout = make([][]rune, len(lines))
	for i, line := range lines {
		DayLayout[i] = []rune(line)
	}

	weatherCodes := map[int]weatherCode{}
	if err := yaml.Unmarshal(weatherCodeBytes, &weatherCodes); err != nil {
		panic(err)
	}

	WeatherCodes = map[int]WeatherCode{}
	for i, code := range weatherCodes {
		sym := strings.Replace(code.Symbol, "\r\n", "\n", -1)
		col := strings.Replace(code.Color, "\r\n", "\n", -1)
		lines := strings.Split(sym, "\n")
		colLines := strings.Split(col, "\n")
		newCode := WeatherCode{
			Name:   code.Name,
			Symbol: make([][]rune, len(lines)),
			Colors: make([][]uint8, len(lines)),
		}
		for i, line := range lines {
			runes := []rune(line)
			runes = runes[1 : len(runes)-1]
			newCode.Symbol[i] = runes
		}
		for i, line := range colLines {
			runes := []rune(line)
			runes = runes[1 : len(runes)-1]
			ids := make([]uint8, len(runes))
			for i, r := range runes {
				id, ok := ColorIDs[r]
				if !ok {
					log.Fatalf("not a color for weather code symbols: '%s'", string(r))
				}
				ids[i] = id
			}
			newCode.Colors[i] = ids
		}

		WeatherCodes[i] = newCode
	}
}
