package filemgr

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/config"
	"io"
	"io/ioutil"
)

type jwtPayload struct {
	Allow []string
}
type jwtSecret struct {
	key   []byte
	token []byte
}

func randSecret() (*jwtSecret, error) {
	sk, err := ioutil.ReadAll(io.LimitReader(rand.Reader, 32))
	if err != nil {
		return nil, err
	}
	p := jwtPayload{
		Allow: api.AllPermissions,
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

/*func (js *jwtSecret) alg() *api.APIAlg {
	return (*api.APIAlg)(jwt.NewHS256(js.key))
}*/

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
