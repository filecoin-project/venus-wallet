package config

// full config
type Config struct {
	API            *APIConfig            `json:"API"`
	DB             *DBConfig             `json:"DB" binding:"required"`
	Metrics        *MetricsConfig        `json:"METRICS"`
	JWT            *JWTConfig            `json:"JWT"`
	Factor         *CryptoFactor         `json:"FACTOR"`
	SignFilter     *SignFilter           `json:"SignFilter"`
	APIRegisterHub *APIRegisterHubConfig `json:"WalletEvent"`
	SignRecorder   *SignRecorderConfig   `json:"SignRecorder"`
}

type APIRegisterHubConfig struct {
	RegisterAPI     []string `json:"apiRegisterHub"`
	Token           string   `json:"token"`
	SupportAccounts []string `json:"supportAccounts"`
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

// aes
type CryptoFactor struct {
	// ScryptN is the N parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	ScryptN int `json:"scryptN"`
	// ScryptP is the P parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	ScryptP int `json:"scryptP"`
}

type SignFilter struct {
	Expr string `json:"expr"`
}

type SignRecorderConfig struct {
	Enable       bool   `json:"enable"`
	KeepDuration string `json:"holdTime"`
}
