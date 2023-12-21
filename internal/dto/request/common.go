package request

type IPInfo struct {
	IP           string
	CountryShort string
	// CountryLong  string
	// Region       string
	// City         string
	// Isp          string
}

type ClientInfo struct {
	AppName    string `json:"app_name" form:"app_name" binding:"required"`
	AppVersion string `json:"app_version" form:"app_version"`
}
