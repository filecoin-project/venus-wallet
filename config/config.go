package config

// full config
type Config struct {
	API     *APIConfig     `json:"API"`
	DB      *DBConfig      `json:"DB" binding:"required"`
	Metrics *MetricsConfig `json:"METRICS"`
	JWT     *JWTConfig     `json:"JWT"`
}

// for keystore
type DBConfig struct {
	Conn      string `json:"conn" binding:"required"`
	Type      string `json:"type" binding:"required"`
	DebugMode bool   `json:"debugMode" binding:"required"`
}

// rpc server address listen
type APIConfig struct {
	ListenAddress string `json:"ListenAddress"`
}

// metrics
type MetricsConfig struct {
	Nickname   string `json:"NickName"`
	HeadNotify bool   `json:"HeadNotify"`
}

// jwt hex token and secret
type JWTConfig struct {
	Token  string `json:"Token"`
	Secret string `json:"Secret"`
}
