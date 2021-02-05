package modules

import (
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"

	"github.com/gbrlsnchs/jwt/v3"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"

	"github.com/ipfs-force-community/venus-wallet/api/apistruct"
	"github.com/ipfs-force-community/venus-wallet/node/modules/dtypes"
	"github.com/ipfs-force-community/venus-wallet/node/repo"
	"github.com/filecoin-project/lotus/chain/types"
)

var log = logging.Logger("modules")


const JWTSecretName = "auth-jwt-private"

type jwtPayload struct {
	Allow []string
}

func APISecret(keystore types.KeyStore, lr repo.LockedRepo) (*dtypes.APIAlg, error) {
	key, err := keystore.Get(JWTSecretName)

	if errors.Is(err, types.ErrKeyInfoNotFound) {
		log.Warn("Generating new API secret")

		sk, err := ioutil.ReadAll(io.LimitReader(rand.Reader, 32))
		if err != nil {
			return nil, err
		}

		key = types.KeyInfo{
			Type:       "jwt-hmac-secret",
			PrivateKey: sk,
		}

		if err := keystore.Put(JWTSecretName, key); err != nil {
			return nil, xerrors.Errorf("writing API secret: %w", err)
		}

		// TODO: make this configurable
		p := jwtPayload{
			Allow: apistruct.AllPermissions,
		}

		cliToken, err := jwt.Sign(&p, jwt.NewHS256(key.PrivateKey))
		if err != nil {
			return nil, err
		}

		if err := lr.SetAPIToken(cliToken); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, xerrors.Errorf("could not get JWT Token: %w", err)
	}

	return (*dtypes.APIAlg)(jwt.NewHS256(key.PrivateKey)), nil
}
