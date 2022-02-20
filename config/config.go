package config

//SqlConnDate 数据库连接模型
type SqlConnDate struct {
	DriverName string `json:"driverName"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Database   string `json:"database"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Charset    string `json:"charset"`
	Loc        string `json:"loc"`
}
