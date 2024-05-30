package api

import (
	"github.com/mlange-42/tom/data"
	"gopkg.in/yaml.v3"
)

var codes map[int]WeatherCode

type WeatherCode struct {
	Name   string
	Symbol string
}

func init() {
	codes = map[int]WeatherCode{}
	if err := yaml.Unmarshal(data.WeatherCodes, &codes); err != nil {
		panic(err)
	}
}

/*var codeText = map[int]string{
	0:  "clear sky",
	1:  "mainly clear",
	2:  "partly cloudy",
	3:  "overcast",
	45: "fog",
	48: "depositing rime fog",
	51: "light drizzle",
	53: "moderate drizzle",
	55: "dense drizzle",
	56: "light freezing drizzle",
	57: "dense freezing drizzle",
	61: "slight rain",
	63: "moderate rain",
	65: "heavy rain",
	66: "light freezing rain",
	67: "heavy freezing rain",
	71: "light snowfall",
	73: "moderate snowfall",
	75: "heavy snowfall",
	77: "snow grains",
	80: "slight rain showers",
	81: "moderate rain showers",
	82: "violent rain showers",
	85: "slight snow showers",
	86: "heavy snow showers",
	95: "slight thunderstorm",
	96: "slight thunderstorm w/ hail",
	97: "heavy thunderstorm",
	99: "heavy thunderstorm w/ hail",
}*/

func WeatherCodeText(code int) (string, bool) {
	s, ok := codes[code]
	return s.Name, ok
}

func WeatherCodeSymbol(code int) (string, bool) {
	s, ok := codes[code]
	return s.Symbol, ok
}
