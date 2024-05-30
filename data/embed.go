package data

import (
	_ "embed"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed weather_codes.yml
var weatherCodes []byte

//go:embed layout.txt
var layout string

var Layout [][]rune

var WeatherCodes map[int]WeatherCode

type WeatherCode struct {
	Name        string
	Symbol      string
	SymbolRunes [][]rune
}

func init() {
	l := strings.Replace(layout, "\r\n", "\n", -1)
	lines := strings.Split(l, "\n")

	Layout = make([][]rune, len(lines))
	for i, line := range lines {
		Layout[i] = []rune(line)
	}

	WeatherCodes = map[int]WeatherCode{}
	if err := yaml.Unmarshal(weatherCodes, &WeatherCodes); err != nil {
		panic(err)
	}

	for i, code := range WeatherCodes {
		sym := strings.Replace(code.Symbol, "\r\n", "\n", -1)
		lines := strings.Split(sym, "\n")
		code.SymbolRunes = make([][]rune, len(lines))
		for i, line := range lines {
			runes := []rune(line)
			runes = runes[1 : len(runes)-1]
			code.SymbolRunes[i] = runes
		}

		WeatherCodes[i] = code
	}
}
