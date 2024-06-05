package config

var DefaultCurrentMetrics = []CurrentMetric{
	CurrentWeatherCode,
	CurrentTemp,
	CurrentApparentTemp,
	CurrentPrecip,
	CurrentCloudCover,
	CurrentWindSpeed,
	CurrentWindDir,
	CurrentWindGusts,
	CurrentRH,
}

var DefaultHourlyMetrics = []HourlyMetric{
	HourlyWeatherCode,
	HourlyTemp,
	HourlyApparentTemp,
	HourlyPrecip,
	HourlyPrecipProb,
	HourlyRH,
	HourlyCloudCover,
	HourlyWindSpeed,
	HourlyWindDir,
	HourlyWindGusts,
}

var DefaultDailyMetrics = []DailyMetric{
	DailyWeatherCode,
	DailyMinTemp,
	DailyMaxTemp,
	DailyPrecip,
	DailyPrecipProb,
	DailySunshine,
	DailyWindSpeed,
	DailyWindDir,
	DailyWindGusts,
}
