package api

import (
	"fmt"
	"strings"
)

type Forecaster string

var (
	Forecast    Forecaster = "forecast"
	DWD                    = "dwd-icon"
	NOAA_US                = "gfs"
	MeteoFrance            = "meteofrance"
	ESMWF                  = "esmwf"
	// TODO: add the other forecasters
)

type Options struct {
	Forecaster        Forecaster // Default "forecast"
	TemperatureUnit   string     // Default "celsius"
	WindSpeedUnit     string     // Default "kmh",
	PrecipitationUnit string     // Default "mm"
	Timezone          string     // Default "UTC"
	Days              int        // Default 7
	PastDays          int        // Default 0
	HourlyMetrics     []string   // Lists required hourly metrics, see https://open-meteo.com/en/docs for valid metrics
	DailyMetrics      []string   // Lists required daily metrics, see https://open-meteo.com/en/docs for valid metrics
}

func (o *Options) ToURL(baseURL string, loc Location) string {
	forecaster := "forecast"
	if o != nil && o.Forecaster != "" {
		forecaster = string(o.Forecaster)
	}
	url := fmt.Sprintf(`%s/%s?latitude=%f&longitude=%f&current_weather=true`, baseURL, forecaster, loc.Lat, loc.Lon)
	if o == nil {
		return url
	}

	if o.TemperatureUnit != "" {
		url = fmt.Sprintf(`%s&temperature_unit=%s`, url, o.TemperatureUnit)
	}
	if o.WindSpeedUnit != "" {
		url = fmt.Sprintf(`%s&windspeed_unit=%s`, url, o.WindSpeedUnit)
	}
	if o.PrecipitationUnit != "" {
		url = fmt.Sprintf(`%s&precipitation_unit=%s`, url, o.PrecipitationUnit)
	}
	if o.Timezone != "" {
		url = fmt.Sprintf(`%s&timezone=%s`, url, o.Timezone)
	}
	if o.Days != 0 {
		url = fmt.Sprintf(`%s&forecast_days=%d`, url, o.Days)
	}
	if o.PastDays != 0 {
		url = fmt.Sprintf(`%s&past_days=%d`, url, o.PastDays)
	}

	if o.HourlyMetrics != nil && len(o.HourlyMetrics) > 0 {
		metrics := strings.Join(o.HourlyMetrics, ",")
		url = fmt.Sprintf(`%s&hourly=%s`, url, metrics)
	}

	if o.DailyMetrics != nil && len(o.DailyMetrics) > 0 {
		metrics := strings.Join(o.DailyMetrics, ",")
		url = fmt.Sprintf(`%s&daily=%s`, url, metrics)
	}

	return url
}
