package filemgr

import (
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/core"
	"os"
	"path/filepath"
)

func (fsr *FsRepo) defConfig() *config.Config {
	return &config.Config{
		API: &config.APIConfig{
			ListenAddress: "/ip4/0.0.0.0/tcp/5678/http",
		},
		DB: &config.DBConfig{
			Conn:      filepath.Join(fsr.path, skKeyStore),
			Type:      "sqlite",
			DebugMode: true,
		},
	}
}

// for program start config init and default element cover
func (fsr *FsRepo) checkConfig(op *OverrideParams) error {
	var (
		cnf *config.Config
		err error
	)
	exist, err := fsr.configExist()
	if err != nil {
		return err
	}
	if exist {
		cnf, err = fsr.loadConfig()
		if err != nil {
			return err
		}
	} else {
		cnf = new(config.Config)
	}
	def := fsr.defConfig()
	reset := false
	if cnf.DB == nil || cnf.DB.Conn == "" {
		cnf.DB = def.DB
		reset = true
	}
	if cnf.API == nil || cnf.API.ListenAddress == "" {
		cnf.API = def.API
		reset = true
	}
	if cnf.JWT == nil || cnf.JWT.Secret == "" {
		cnf.JWT, err = RandJWTConfig()
		if err != nil {
			return err
		}
		reset = true
	}
	if cnf.Factor == nil || cnf.Factor.ScryptP == 0 && cnf.Factor.ScryptN == 0 {
		cnf.Factor = &config.CryptoFactor{
			ScryptN: 1 << 18,
			ScryptP: 1,
		}
	}
	if cnf.Strategy == nil {
		cnf.Strategy = &config.StrategyConfig{
			Level:   0,
			NodeURL: "",
		}
	}
	if op != nil {
		// override
		if op.API != core.StringEmpty {
			cnf.API.ListenAddress = op.API
		}
	}
	if reset {
		err = config.CoverConfig(fsr.configPath(), cnf)
		if err != nil {
			return err
		}
	}
	fsr.cnf = cnf
	return nil
}

func (fsr *FsRepo) configExist() (bool, error) {
	_, err := os.Stat(fsr.configPath())
	if err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}
func (fsr *FsRepo) loadConfig() (*config.Config, error) {
	cnf, err := config.DecodeConfig(fsr.configPath())
	if err != nil {
		return nil, err
	}
	return cnf, nil
}
func (fsr *FsRepo) configPath() string {
	return filepath.Join(fsr.path, skConfig)
}
