package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/mlange-42/tom/util/agg"
)

type Forecaster string

const (
	Forecast    Forecaster = "forecast"
	DWD         Forecaster = "dwd-icon"
	NOAA_US     Forecaster = "gfs"
	MeteoFrance Forecaster = "meteofrance"
	ESMWF       Forecaster = "esmwf"
	// TODO: add the other forecasters
)

type CurrentMetric string

const (
	CurrentWeatherCode  CurrentMetric = "weather_code"
	CurrentTemp         CurrentMetric = "temperature_2m"
	CurrentApparentTemp CurrentMetric = "apparent_temperature"
	CurrentRH           CurrentMetric = "relative_humidity_2m"
	CurrentPrecip       CurrentMetric = "precipitation"
	CurrentCloudCover   CurrentMetric = "cloud_cover"
	CurrentWindSpeed    CurrentMetric = "wind_speed_10m"
	CurrentWindDir      CurrentMetric = "wind_direction_10m"
)

type HourlyMetric string

const (
	HourlyWeatherCode  HourlyMetric = "weather_code"
	HourlyTemp         HourlyMetric = "temperature_2m"
	HourlyApparentTemp HourlyMetric = "apparent_temperature"
	HourlyRH           HourlyMetric = "relative_humidity_2m"
	HourlyPrecipProb   HourlyMetric = "precipitation_probability"
	HourlyPrecip       HourlyMetric = "precipitation"
	HourlyCloudCover   HourlyMetric = "cloud_cover"
	HourlyWindSpeed    HourlyMetric = "wind_speed_10m"
	HourlyWindDir      HourlyMetric = "wind_direction_10m"
	HourlyWindGusts    HourlyMetric = "wind_gusts_10m"
)

type DailyMetric string

const (
	DailyWeatherCode DailyMetric = "weather_code"
	DailyMinTemp     DailyMetric = "temperature_2m_min"
	DailyMaxTemp     DailyMetric = "temperature_2m_max"
	DailySunshine    DailyMetric = "sunshine_duration"
	DailyPrecip      DailyMetric = "precipitation_sum"
	DailyPrecipHours DailyMetric = "precipitation_hours"
	DailyPrecipProb  DailyMetric = "precipitation_probability_max"
	DailyWindSpeed   DailyMetric = "wind_speed_10m_max"
	DailyWindGusts   DailyMetric = "wind_gusts_10m_max"
	DailyWindDir     DailyMetric = "wind_direction_10m_dominant"
)

var aggregators = map[HourlyMetric]agg.Aggregator{
	HourlyWeatherCode:  &agg.Point{},
	HourlyTemp:         &agg.Point{},
	HourlyApparentTemp: &agg.Point{},
	HourlyRH:           &agg.Point{},
	HourlyPrecipProb:   &agg.Max{},
	HourlyPrecip:       &agg.Sum{},
	HourlyCloudCover:   &agg.Point{},
	HourlyWindSpeed:    &agg.Max{},
	HourlyWindGusts:    &agg.Max{},
	HourlyWindDir:      &agg.Point{},
}

type Options interface {
	ToURL(baseURL string) string
}

type ForecastOptions struct {
	Location          Location
	Forecaster        Forecaster      // Default "forecast"
	TemperatureUnit   string          // Default "celsius"
	WindSpeedUnit     string          // Default "kmh",
	PrecipitationUnit string          // Default "mm"
	Days              int             // Default 7
	PastDays          int             // Default 0
	CurrentMetrics    []CurrentMetric // List of required current metrics
	HourlyMetrics     []HourlyMetric  // List of required hourly metrics
	DailyMetrics      []DailyMetric   // List of required daily metrics
}

func (o *ForecastOptions) ToURL(baseURL string) string {
	forecaster := "forecast"
	if o.Forecaster != "" {
		forecaster = string(o.Forecaster)
	}
	url := fmt.Sprintf(`%s/%s?latitude=%f&longitude=%f`, baseURL, forecaster, o.Location.Lat, o.Location.Lon)

	if o.TemperatureUnit != "" {
		url = fmt.Sprintf(`%s&temperature_unit=%s`, url, o.TemperatureUnit)
	}
	if o.WindSpeedUnit != "" {
		url = fmt.Sprintf(`%s&windspeed_unit=%s`, url, o.WindSpeedUnit)
	}
	if o.PrecipitationUnit != "" {
		url = fmt.Sprintf(`%s&precipitation_unit=%s`, url, o.PrecipitationUnit)
	}
	if o.Location.TimeZone != "" {
		url = fmt.Sprintf(`%s&timezone=%s`, url, o.Location.TimeZone)
	}
	if o.Days != 0 {
		url = fmt.Sprintf(`%s&forecast_days=%d`, url, o.Days)
	}
	if o.PastDays != 0 {
		url = fmt.Sprintf(`%s&past_days=%d`, url, o.PastDays)
	}

	if len(o.CurrentMetrics) > 0 {
		met := make([]string, len(o.CurrentMetrics))
		for i, m := range o.CurrentMetrics {
			met[i] = string(m)
		}
		metrics := strings.Join(met, ",")
		url = fmt.Sprintf(`%s&current=%s`, url, metrics)
	}

	if len(o.HourlyMetrics) > 0 {
		met := make([]string, len(o.HourlyMetrics))
		for i, m := range o.HourlyMetrics {
			met[i] = string(m)
		}
		metrics := strings.Join(met, ",")
		url = fmt.Sprintf(`%s&hourly=%s`, url, metrics)
	}

	if len(o.DailyMetrics) > 0 {
		met := make([]string, len(o.DailyMetrics))
		for i, m := range o.DailyMetrics {
			met[i] = string(m)
		}
		metrics := strings.Join(met, ",")
		url = fmt.Sprintf(`%s&daily=%s`, url, metrics)
	}

	return url
}

type GeoOptions struct {
	Name  string
	Count int
}

func (o *GeoOptions) ToURL(baseURL string) string {
	url := fmt.Sprintf("%s/search?name=%s", baseURL, url.QueryEscape(o.Name))

	if o.Count != 0 {
		url = fmt.Sprintf("%s&count=%d", url, o.Count)
	}

	return url
}
