package config

type Service struct {
	Name        string
	Description string
	UrlName     string
}

var Services = []Service{
	{
		Name:        "OM",
		Description: "Open-Meteo",
		UrlName:     "forecast",
	},
	{
		Name:        "DWD",
		Description: "DWD Germany",
		UrlName:     "dwd-icon",
	},
	{
		Name:        "NOAA",
		Description: "NOAA U.S.",
		UrlName:     "gfs",
	},
	{
		Name:        "MF",
		Description: "Meteo France",
		UrlName:     "meteofrance",
	},
	{
		Name:        "ECMWF",
		Description: "ECMWF",
		UrlName:     "esmwf",
	},
	// TODO: add the other available services
}
