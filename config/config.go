package config

// full config
type Config struct {
	API     *APIConfig     `json:"API"`
	DB      *DBConfig      `json:"DB" binding:"required"`
	Metrics *MetricsConfig `json:"METRICS"`
	JWT     *JWTConfig     `json:"JWT"`
	Factor  *CryptoFactor  `json:"FACTOR"`
}

// for keystore
type DBConfig struct {
	Conn      string `json:"conn" binding:"required"`
	Type      string `json:"type" binding:"required"`
	DebugMode bool   `json:"debugMode" binding:"required"`
}

// rpc server address listen
type APIConfig struct {
	ListenAddress string `json:"listenAddress"`
}

// metrics
type MetricsConfig struct {
	Nickname   string `json:"nickName"`
	HeadNotify bool   `json:"headNotify"`
}

// jwt hex token and secret
type JWTConfig struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}

type CryptoFactor struct {
	// ScryptN is the N parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	ScryptN int `json:"scryptN"`
	// ScryptP is the P parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	ScryptP int `json:"scryptP"`
}
