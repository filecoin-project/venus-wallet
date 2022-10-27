package filemgr

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/crypto/aes"
	jwt "github.com/gbrlsnchs/jwt/v3"
	"github.com/google/uuid"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	ErrNoAPIEndpoint     = errors.New("API not running (no endpoint)")
	ErrNoAPIToken        = errors.New("API token not set")
	ErrRepoAlreadyLocked = errors.New("repo is already locked")
	ErrClosedRepo        = errors.New("repo is no longer open")
)

// FsRepo is struct for repo, use NewFS to create
type FsRepo struct {
	path string
	cnf  *config.Config
}

type OverrideParams struct {
	API string

	GatewayAPI      []string
	GatewayToken    string
	SupportAccounts []string
}

var _ Repo = &FsRepo{}

// NewFS creates a repo instance based on a path on file system
func NewFS(path string, op *OverrideParams) (Repo, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}
	fs := &FsRepo{
		path: path,
	}
	err = fs.init()
	if err != nil {
		return nil, err
	}
	err = fs.checkConfig(op)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func (fsr *FsRepo) APISecret() (*jwt.HMACSHA, error) {
	sec, err := hex.DecodeString(fsr.cnf.JWT.Secret)
	if err != nil {
		return nil, err
	}
	return jwt.NewHS256(sec), nil
}

func (fsr *FsRepo) init() error {
	exist, err := fsr.exists()
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	err = os.Mkdir(fsr.path, 0o755) //nolint: gosec
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func (fsr *FsRepo) exists() (bool, error) {
	_, err := os.Stat(filepath.Join(fsr.path, skKeyStore))
	notexist := os.IsNotExist(err)
	if notexist {
		err = nil
	}
	return !notexist, err
}

func (fsr *FsRepo) Config() *config.Config {
	return fsr.cnf
}

// APIEndpoint returns endpoint of API in this repo
func (fsr *FsRepo) APIEndpoint() (string, error) {
	strma := strings.TrimSpace(fsr.cnf.API.ListenAddress)
	return strma, nil
}

func (fsr *FsRepo) APIToken() ([]byte, error) {
	return hex.DecodeString(fsr.cnf.JWT.Token)
}

func (fsr *FsRepo) APIStrategyToken(password string) (string, error) {
	hashPasswd := aes.Keccak256([]byte(password))
	rootKey, err := aes.EncryptData(hashPasswd, []byte("root"), fsr.cnf.Factor.ScryptN, fsr.cnf.Factor.ScryptP)
	if err != nil {
		return "", errors.New("failed to gen token seed")
	}
	rootKB, err := json.Marshal(rootKey)
	if err != nil {
		return "", errors.New("failed to marshal token seed")
	}
	rootk, err := uuid.NewRandomFromReader(bytes.NewBuffer(rootKB))
	if err != nil {
		return "", errors.New("failed to convert token seed to uuid")
	}
	return rootk.String(), nil
}
