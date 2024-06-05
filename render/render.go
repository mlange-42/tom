package render

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mlange-42/tom/config"
	"github.com/mlange-42/tom/data"
	"github.com/mlange-42/tom/util"
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
	layout := make([][]rune, len(data.DayLayout))
	colors := make([][]uint8, len(data.DayLayout))
	for i, runes := range data.DayLayout {
		layout[i] = append(layout[i], runes...)
		colors[i] = make([]uint8, len(runes))
	}

	layoutWidth := len(layout[0])
	boxWidth := (layoutWidth-1)/4 - 1
	yOffset := 1

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
			copy(colors[j+yOffset+1][x:x+len(line)], codeProps.Colors[j])
		}

		text := []string{
			fmt.Sprintf("%-5s %s", timeStart[idx].Format(config.TimeLayout), codeProps.Name),
			fmt.Sprintf("%2d (%2d) °C", int(math.Round(temp[idx])), int(math.Round(appTemp[idx]))),
			fmt.Sprintf("%4.1fmm/%3d%%", precip[idx], int(math.Round(precipProb[idx]))),
			fmt.Sprintf("%3dkm/h %-2s", int(math.Round(wind[idx])), config.Direction(windDir[idx])),
			fmt.Sprintf("%3d%%C %3d%%H", int(math.Round(clouds[idx])), int(math.Round(humidity[idx]))),
		}
		cols := make([][]uint8, len(text))

		t, _, p, w := calcColors(temp[idx], 0, precip[idx], precipProb[idx], wind[idx])
		cols[1] = util.Repeat(t, utf8.RuneCountInString(text[1]))
		cols[2] = util.Repeat(p, utf8.RuneCountInString(text[2]))
		cols[3] = util.Repeat(w, utf8.RuneCountInString(text[3]))

		symWidth := len(codeProps.Symbol[0])
		x += 1
		for j, line := range text {
			maxLen := boxWidth - (symWidth + 1)
			if j == 0 {
				maxLen = boxWidth - 1
			}
			len := MinInt(maxLen, len(line))
			copy(layout[j+yOffset][x:x+len], []rune(line[:len]))
			copy(colors[j+yOffset][x:x+len], cols[j])

			if j == 0 {
				x += symWidth
			}
		}
	}

	result := make([]string, len(layout))
	builder := strings.Builder{}

	var prevColor uint8
	for i, runes := range layout {
		cols := colors[i]

		for j, r := range runes {
			c := cols[j]
			if c != prevColor {
				builder.WriteString(config.Colors[c].Tag)
			}
			builder.WriteRune(r)
			prevColor = c
		}

		builder.WriteString(config.Colors[0].Tag)
		result[i] = builder.String()
		builder.Reset()
	}

	resultStr := strings.Join(result, "\n")
	return resultStr
}

func (r *Renderer) DaySummary(index int) string {
	code := int(r.data.GetDaily(config.DailyWeatherCode)[index])
	codeProps, ok := data.WeatherCodes[code]
	if !ok {
		log.Fatalf("unknown weather code %d", code)
	}

	minTemp := r.data.GetDaily(config.DailyMinTemp)[index]
	maxTemp := r.data.GetDaily(config.DailyMaxTemp)[index]
	precip := r.data.GetDaily(config.DailyPrecip)[index]
	precipProb := r.data.GetDaily(config.DailyPrecipProb)[index]
	windSpeed := r.data.GetDaily(config.DailyWindSpeed)[index]

	t1, t2, p, w := calcColors(minTemp, maxTemp, precip, precipProb, windSpeed)

	return fmt.Sprintf(
		"%-27s %s%2d[-]-%s%2d[-]°C  %s%4.1fmm/%3d%%  %s%3dkm/h %-2s[-]",
		codeProps.Name,
		config.Colors[t1].Tag,
		int(math.Round(minTemp)),
		config.Colors[t2].Tag,
		int(math.Round(maxTemp)),
		config.Colors[p].Tag,
		precip,
		int(math.Round(precipProb)),
		config.Colors[w].Tag,
		int(windSpeed),
		config.Direction(r.data.GetDaily(config.DailyWindDir)[index]),
	)
}

func (r *Renderer) Current() string {
	code := int(r.data.GetCurrent(config.CurrentWeatherCode))
	codeProps, ok := data.WeatherCodes[code]
	if !ok {
		log.Fatalf("unknown weather code %d", code)
	}

	temp := r.data.GetCurrent(config.CurrentTemp)
	precip := r.data.GetCurrent(config.CurrentPrecip)
	windSpeed := r.data.GetCurrent(config.CurrentWindSpeed)

	t, _, p, w := calcColors(temp, 0, precip, 0, windSpeed)

	return fmt.Sprintf(
		"%s %s%3d°C  %s%4.1fmm/h  %s%3dkm/h %-2s[-]  %3d%%C  %3d%%H",
		codeProps.Name,
		config.Colors[t].Tag,
		int(math.Round(temp)),
		config.Colors[p].Tag,
		precip,
		config.Colors[w].Tag,
		int(windSpeed),
		config.Direction(r.data.GetCurrent(config.CurrentWindDir)),
		int(r.data.GetCurrent(config.CurrentCloudCover)),
		int(r.data.GetCurrent(config.CurrentRH)),
	)
}

func (r *Renderer) Charts(now time.Time) string {
	builder := strings.Builder{}
	builder.WriteString("Temperature [°C]\n")
	builder.WriteString(r.chart(config.HourlyTemp, false, now) + "\n")
	builder.WriteString("\nPrecipitation [mm/h]\n")
	builder.WriteString(r.chart(config.HourlyPrecip, true, now) + "\n")
	builder.WriteString("\nPrecipitation probability [%]\n")
	builder.WriteString(r.chart(config.HourlyPrecipProb, true, now) + "\n")
	builder.WriteString("\nWind speed [km/h]\n")
	builder.WriteString(r.chart(config.HourlyWindSpeed, false, now) + "\n")
	builder.WriteString("\nCloud cover [%]\n")
	builder.WriteString(r.chart(config.HourlyCloudCover, true, now) + "\n")
	builder.WriteString("\nRelative humidity [%]\n")
	builder.WriteString(r.chart(config.HourlyRH, true, now))

	return builder.String()
}

func (r *Renderer) chart(metric config.HourlyMetric, bars bool, now time.Time) string {
	chart := NewChart(len(r.data.HourlyTime)/2, 6)
	vMin, vMax := chart.Series(r.data.GetHourly(metric), bars)

	for i := 0; i < len(r.data.HourlyTime); i += 24 {
		chart.VLine(i, 4, bars)
	}

	runes := chart.Runes()

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%8.1f ", vMax))
	builder.WriteString(fmt.Sprintf("%s\n", string(runes[0])))

	for i := 1; i < len(runes)-1; i++ {
		builder.WriteString(fmt.Sprintf("         %s\n", string(runes[i])))
	}

	builder.WriteString(fmt.Sprintf("%8.1f ", vMin))
	builder.WriteString(fmt.Sprintf("%s\n", string(runes[len(runes)-1])))

	builder.WriteString("         ")
	for _, t := range r.data.DailyTime {
		ts := t.Format(config.DateLayoutShort)
		if now.YearDay() == t.YearDay() {
			builder.WriteString(fmt.Sprintf("[yellow]%11s[-] ", ts))
		} else {
			builder.WriteString(fmt.Sprintf("%11s ", ts))
		}
	}

	return builder.String()
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func calcColors(temp1, temp2, precip, precipProb, wind float64) (t1 uint8, t2 uint8, p uint8, w uint8) {
	if temp1 <= 0 {
		t1 = 3 // blue
	} else if temp1 > 30 {
		t1 = 2 // red
	} else if temp1 > 20 {
		t1 = 1 // yellow
	}
	if temp2 <= 0 {
		t2 = 3 // blue
	} else if temp2 > 30 {
		t2 = 2 // red
	} else if temp2 > 20 {
		t2 = 1 // yellow
	}

	if precip >= 1 || precipProb > 50 {
		p = 3 // blue
	}

	if wind >= 62 {
		w = 2 // red
	} else if wind >= 29 {
		w = 1 // yellow
	}

	return
}
