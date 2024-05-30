package app

import (
	"fmt"
	"log"
	"math"

	"github.com/mlange-42/tom/api"
	"github.com/mlange-42/tom/data"
)

func formatCurrent(d *api.MeteoResult) string {
	code := int(d.GetCurrent(api.CurrentWeatherCode))
	codeProps, ok := data.WeatherCodes[code]
	if !ok {
		log.Fatalf("unknown weather code %d", code)
	}

	return fmt.Sprintf(
		"%s %3d°C  %4.1fmm  %3dkm/h %-2s  %3d%%CC  %3d%%RH",
		codeProps.Name, int(math.Round(d.GetCurrent(api.CurrentTemp))),
		d.GetCurrent(api.CurrentPrecip),
		int(d.GetCurrent(api.CurrentWindSpeed)),
		api.Direction(d.GetCurrent(api.CurrentWindDir)),
		int(d.GetCurrent(api.CurrentCloudCover)),
		int(d.GetCurrent(api.CurrentRH)),
	)
}
