package common

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/build"
	"github.com/ipfs-force-community/venus-wallet/node/modules/dtypes"
	"github.com/gbrlsnchs/jwt/v3"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

var log = logging.Logger("common")

type CommonAPI struct {
	fx.In
	APISecret *dtypes.APIAlg
}

type jwtPayload struct {
	Allow []string
}

func (a *CommonAPI) AuthVerify(ctx context.Context, token string) ([]api.Permission, error) {
	var payload jwtPayload
	if _, err := jwt.Verify([]byte(token), (*jwt.HMACSHA)(a.APISecret), &payload); err != nil {
		return nil, xerrors.Errorf("JWT Verification failed: %w", err)
	}

	return payload.Allow, nil
}

func (a *CommonAPI) AuthNew(ctx context.Context, perms []api.Permission) ([]byte, error) {
	p := jwtPayload{
		Allow: perms, // TODO: consider checking validity
	}

	return jwt.Sign(&p, (*jwt.HMACSHA)(a.APISecret))
}

func (a *CommonAPI) Version(context.Context) (api.Version, error) {
	return api.Version{
		Version:    build.UserVersion,
		APIVersion: build.APIVersion,
	}, nil
}

func (a *CommonAPI) LogList(context.Context) ([]string, error) {
	return logging.GetSubsystems(), nil
}

func (a *CommonAPI) LogSetLevel(ctx context.Context, subsystem, level string) error {
	return logging.SetLogLevel(subsystem, level)
}

var _ api.Common = &CommonAPI{}
