package common

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/filecoin-project/venus-wallet/version"
	api "github.com/filecoin-project/venus/venus-shared/api/wallet"
	"github.com/filecoin-project/venus/venus-shared/types"
	"github.com/gbrlsnchs/jwt/v3"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
)

type ICommon = api.ICommon

type APIAlg jwt.HMACSHA

var _ api.ICommon = &Common{}

type Common struct {
	fx.In
	APISecret *APIAlg
}

type jwtPayload struct {
	Allow []string
}

func (a *Common) AuthVerify(ctx context.Context, token string) ([]auth.Permission, error) {
	var payload jwtPayload
	if _, err := jwt.Verify([]byte(token), (*jwt.HMACSHA)(a.APISecret), &payload); err != nil {
		return nil, fmt.Errorf("JWT Verification failed: %w", err)
	}
	return payload.Allow, nil
}

func (a *Common) AuthNew(ctx context.Context, perms []auth.Permission) ([]byte, error) {
	p := jwtPayload{
		Allow: perms, // TODO: consider checking validity
	}
	return jwt.Sign(&p, (*jwt.HMACSHA)(a.APISecret))
}

func (a *Common) Version(context.Context) (types.Version, error) {
	return types.Version{
		Version:    version.UserVersion,
		APIVersion: version.APIVersion,
	}, nil
}

func (a *Common) LogList(context.Context) ([]string, error) {
	return logging.GetSubsystems(), nil
}

func (a *Common) LogSetLevel(ctx context.Context, subsystem, level string) error {
	return logging.SetLogLevel(subsystem, level)
}
