package render

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/mlange-42/tom/config"
	"github.com/mlange-42/tom/data"
)

type Renderer struct {
	data *config.MeteoResult
}

func NewRenderer(data *config.MeteoResult) Renderer {
	return Renderer{
		data: data,
	}
}

func (r *Renderer) DaySixHourly(index int) string {
	layout := make([][]rune, len(data.Layout))
	colors := make([][]uint8, len(data.Layout))
	for i, runes := range data.Layout {
		layout[i] = append(layout[i], runes...)
		colors[i] = make([]uint8, len(runes))
	}

	layoutWidth := len(layout[0])
	boxWidth := (layoutWidth-1)/4 - 1
	yOffset := 1

	_ = boxWidth

	timeStart := r.data.SixHourlyTime
	codes := r.data.GetSixHourly(config.HourlyWeatherCode)
	temp := r.data.GetSixHourly(config.HourlyTemp)
	appTemp := r.data.GetSixHourly(config.HourlyApparentTemp)
	precip := r.data.GetSixHourly(config.HourlyPrecip)
	precipProb := r.data.GetSixHourly(config.HourlyPrecipProb)
	wind := r.data.GetSixHourly(config.HourlyWindSpeed)
	windDir := r.data.GetSixHourly(config.HourlyWindDir)
	clouds := r.data.GetSixHourly(config.HourlyCloudCover)
	humidity := r.data.GetSixHourly(config.HourlyRH)

	for i := 0; i < 4; i++ {
		idx := index + i

		code := int(codes[idx])
		codeProps, ok := data.WeatherCodes[code]
		if !ok {
			log.Fatalf("unknown weather code %d", code)
		}

		x := 1 + (boxWidth+1)*i
		for j, line := range codeProps.Symbol {
			copy(layout[j+yOffset+1][x:x+len(line)], line)
		}

		text := []string{
			fmt.Sprintf("%-5s %s", timeStart[idx].Format(config.TimeLayout), codeProps.Name),
			fmt.Sprintf("%2d (%2d) °C", int(math.Round(temp[idx])), int(math.Round(appTemp[idx]))),
			fmt.Sprintf("%4.1fmm/%3d%%", precip[idx], int(math.Round(precipProb[idx]))),
			fmt.Sprintf("%3dkm/h %-2s", int(math.Round(wind[idx])), config.Direction(windDir[idx])),
			fmt.Sprintf("%3d%%CC %3d%%RH", int(math.Round(clouds[idx])), int(math.Round(humidity[idx]))),
		}
		symWidth := len(codeProps.Symbol[0])
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

func (r *Renderer) DaySummary(index int) string {
	code := int(r.data.GetDaily(config.DailyWeatherCode)[index])
	codeProps, ok := data.WeatherCodes[code]
	if !ok {
		log.Fatalf("unknown weather code %d", code)
	}

	return fmt.Sprintf(
		"%-27s %2d-%2d°C  %4.1fmm/%3d%%  %3dkm/h %-2s",
		codeProps.Name,
		int(math.Round(r.data.GetDaily(config.DailyMinTemp)[index])),
		int(math.Round(r.data.GetDaily(config.DailyMaxTemp)[index])),
		r.data.GetDaily(config.DailyPrecip)[index],
		int(math.Round(r.data.GetDaily(config.DailyPrecipProb)[index])),
		int(r.data.GetDaily(config.DailyWindSpeed)[index]),
		config.Direction(r.data.GetDaily(config.DailyWindDir)[index]),
	)
}

func (r *Renderer) Current() string {
	code := int(r.data.GetCurrent(config.CurrentWeatherCode))
	codeProps, ok := data.WeatherCodes[code]
	if !ok {
		log.Fatalf("unknown weather code %d", code)
	}

	return fmt.Sprintf(
		"%s %3d°C  %4.1fmm  %3dkm/h %-2s  %3d%%CC  %3d%%RH",
		codeProps.Name, int(math.Round(r.data.GetCurrent(config.CurrentTemp))),
		r.data.GetCurrent(config.CurrentPrecip),
		int(r.data.GetCurrent(config.CurrentWindSpeed)),
		config.Direction(r.data.GetCurrent(config.CurrentWindDir)),
		int(r.data.GetCurrent(config.CurrentCloudCover)),
		int(r.data.GetCurrent(config.CurrentRH)),
	)
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
