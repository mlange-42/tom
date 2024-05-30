package render

import (
	"fmt"
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

	for i := 0; i < 4; i++ {
		idx := index + i

		code := int(r.data.GetHourly(api.HourlyWeatherCode)[idx])
		codeProps, ok := data.WeatherCodes[code]
		if !ok {
			codeProps = data.WeatherCodes[0]
		}

		x := 1 + (boxWidth+1)*i
		for j, line := range codeProps.SymbolRunes {
			copy(layout[j+yOffset+1][x:x+len(line)], line)
		}

		text := []string{
			fmt.Sprintf("%-5s %s", r.data.SixHourlyTime[idx].Format(api.TimeLayout), codeProps.Name),
		}
		symWidth := len(codeProps.SymbolRunes[0])
		x += 1
		for j, line := range text {
			maxLen := boxWidth - (symWidth + 1)
			if j == 0 {
				maxLen = boxWidth
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
