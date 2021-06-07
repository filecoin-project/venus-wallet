package config

// full config
type Config struct {
	API            *APIConfig            `json:"API"`
	DB             *DBConfig             `json:"DB" binding:"required"`
	Metrics        *MetricsConfig        `json:"METRICS"`
	JWT            *JWTConfig            `json:"JWT"`
	Factor         *CryptoFactor         `json:"FACTOR"`
	Strategy       *StrategyConfig       `json:"STRATEGY"`
	APIRegisterHub *APIRegisterHubConfig `json:"WalletEvent"`
}

// strategy validation
type StrategyConfig struct {
	Level   uint8  `json:"level"`   // 0：nouse  1：only check struct  2：check struct and msg.method
	NodeURL string `json:"nodeUrl"` // need config when Level = 2 and get the actor for msg.to
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
