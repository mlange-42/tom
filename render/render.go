package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/mlange-42/tom/api"
	"github.com/mlange-42/tom/data"
)

type Renderer struct {
	data *api.MeteoResult
}

func NewRenderer(data *api.MeteoResult) Renderer {
	return Renderer{
		data: data,
	}
}

func (r *Renderer) Day(index int) string {
	layout := make([][]rune, len(data.Layout))
	for i, runes := range data.Layout {
		layout[i] = append(layout[i], runes...)
	}
	layoutWidth := len(layout[0])
	boxWidth := (layoutWidth-1)/4 - 1
	yOffset := 1

	_ = boxWidth

	codes := r.data.GetSixHourly(api.HourlyWeatherCode)
	temp := r.data.GetSixHourly(api.HourlyTemp)
	appTemp := r.data.GetSixHourly(api.HourlyApparentTemp)
	precip := r.data.GetSixHourly(api.HourlyPrecip)
	precipProb := r.data.GetSixHourly(api.HourlyPrecipProb)
	wind := r.data.GetSixHourly(api.HourlyWindSpeed)
	windDir := r.data.GetSixHourly(api.HourlyWindDir)
	clouds := r.data.GetSixHourly(api.HourlyCloudCover)
	humidity := r.data.GetSixHourly(api.HourlyRH)

	for i := 0; i < 4; i++ {
		idx := index + i

		code := int(codes[idx])
		codeProps, ok := data.WeatherCodes[code]
		if !ok {
			panic(fmt.Sprintf("unknown weather code %d", code))
			//codeProps = data.WeatherCodes[0]
		}

		x := 1 + (boxWidth+1)*i
		for j, line := range codeProps.SymbolRunes {
			copy(layout[j+yOffset+1][x:x+len(line)], line)
		}

		text := []string{
			fmt.Sprintf("%-5s %s", r.data.SixHourlyTime[idx].Format(api.TimeLayout), codeProps.Name),
			fmt.Sprintf("%2d (%2d) °C", int(math.Round(temp[idx])), int(math.Round(appTemp[idx]))),
			fmt.Sprintf("%2dmm (%3d%%)", int(math.Round(precip[idx])), int(math.Round(precipProb[idx]))),
			fmt.Sprintf("%3dkm/h %-2s", int(math.Round(wind[idx])), api.Direction(windDir[idx])),
			fmt.Sprintf("%3d%%CC %3d%%RH", int(math.Round(clouds[idx])), int(math.Round(humidity[idx]))),
		}
		symWidth := len(codeProps.SymbolRunes[0])
		x += 1
		for j, line := range text {
			maxLen := boxWidth - (symWidth + 1)
			if j == 0 {
				maxLen = boxWidth - 1
			}
			len := MinInt(maxLen, len(line))
			copy(layout[j+yOffset][x:x+len], []rune(line[:len]))

			if j == 0 {
				x += symWidth
			}
		}
	}

	result := make([]string, len(layout))
	for i, runes := range layout {
		result[i] = string(runes)
	}

	return strings.Join(result, "\n")
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
