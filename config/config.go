package config

type SqlConnDate struct {
	DriverName string `json:"driver_name"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Database   string `json:"database"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Charset    string `json:"charset"`
	Loc        string `json:"loc"`
}
