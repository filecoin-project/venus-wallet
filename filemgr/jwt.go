package filemgr

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	jwt "github.com/gbrlsnchs/jwt/v3"

	"github.com/ipfs-force-community/sophon-auth/core"

	"github.com/filecoin-project/venus-wallet/config"
)

type jwtPayload struct {
	Allow []string
}

type jwtSecret struct {
	key   []byte
	token []byte
}

// Random generation of secret keys
func randSecret() (*jwtSecret, error) {
	allPermissions := core.AdaptOldStrategy(core.PermAdmin)
	p := jwtPayload{
		Allow: allPermissions,
	}

	sk, err := io.ReadAll(io.LimitReader(rand.Reader, 32))
	if err != nil {
		return nil, err
	}

	cliToken, err := jwt.Sign(&p, jwt.NewHS256(sk))
	if err != nil {
		return nil, err
	}

	return &jwtSecret{
		key:   sk,
		token: cliToken,
	}, nil
}

// Random generation of JWT config
func RandJWTConfig() (*config.JWTConfig, error) {
	js, err := randSecret()
	if err != nil {
		return nil, err
	}
	cnf := &config.JWTConfig{
		Token:  hex.EncodeToString(js.token),
		Secret: hex.EncodeToString(js.key),
	}
	return cnf, nil
}
