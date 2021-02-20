package config

type Config struct {
	API     *APIConfig     `json:"API"`
	DB      *DBConfig      `json:"DB" binding:"required"`
	Metrics *MetricsConfig `json:"METRICS"`
	JWT     *JWTConfig     `json:"JWT"`
}

type DBConfig struct {
	Conn      string `json:"conn" binding:"required"`
	Type      string `json:"type" binding:"required"`
	DebugMode bool   `json:"debugMode" binding:"required"`
}

type APIConfig struct {
	ListenAddress string `json:"ListenAddress"`
}

type MetricsConfig struct {
	Nickname   string `json:"NickName"`
	HeadNotify bool   `json:"HeadNotify"`
}

type JWTConfig struct {
	Token  string `json:"Token"`
	Secret string `json:"Secret"`
}
