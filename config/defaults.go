package config

var DefaultCurrentMetrics = []CurrentMetric{
	CurrentWeatherCode,
	CurrentTemp,
	CurrentApparentTemp,
	CurrentPrecip,
	CurrentCloudCover,
	CurrentWindSpeed,
	CurrentWindDir,
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
}
