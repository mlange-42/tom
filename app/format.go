package app

import (
	"fmt"
	"log"
	"math"

	"github.com/mlange-42/tom/api"
)

func formatCurrent(data *api.MeteoResult) string {
	code := int(data.GetCurrent(api.CurrentWeatherCode))
	codeStr, ok := api.WeatherCodeText(code)
	if !ok {
		log.Fatalf("unknown weather code %d", code)
	}

	return fmt.Sprintf(
		"%s  %3d°C  %4.1fmm  %3dkm/h %-3s  %3d%%CC  %3d%%RH",
		codeStr, int(math.Round(data.GetCurrent(api.CurrentTemp))),
		data.GetCurrent(api.CurrentPrecip),
		int(data.GetCurrent(api.CurrentWindSpeed)),
		api.Direction(data.GetCurrent(api.CurrentWindDir)),
		int(data.GetCurrent(api.CurrentCloudCover)),
		int(data.GetCurrent(api.CurrentRH)),
	)
}

func formatSixHourly(data *api.MeteoResult, idx int) string {
	code := int(data.GetSixHourly(api.HourlyWeatherCode)[idx])
	codeStr, ok := api.WeatherCodeText(code)
	if !ok {
		log.Fatalf("unknown weather code %d", code)
	}

	return fmt.Sprintf(
		"%27s  %3d°C  %4.1fmm (%3d%%)  %3dkm/h %-3s  %3d%%CC  %3d%%RH",
		codeStr, int(math.Round(data.GetSixHourly(api.HourlyTemp)[idx])),
		data.GetSixHourly(api.HourlyPrecip)[idx],
		int(data.GetSixHourly(api.HourlyPrecipProb)[idx]),
		int(data.GetSixHourly(api.HourlyWindSpeed)[idx]),
		api.Direction(data.GetSixHourly(api.HourlyWindDir)[idx]),
		int(data.GetSixHourly(api.HourlyCloudCover)[idx]),
		int(data.GetSixHourly(api.HourlyRH)[idx]),
	)
}
