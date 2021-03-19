package common

import (
	"context"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/ipfs-force-community/venus-wallet/api/permission"
	"github.com/ipfs-force-community/venus-wallet/version"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

type ICommon interface {
	// Auth
	AuthVerify(ctx context.Context, token string) ([]permission.Permission, error)
	AuthNew(ctx context.Context, perms []permission.Permission) ([]byte, error)

	// Version provides information about API provider
	Version(context.Context) (Version, error)

	LogList(context.Context) ([]string, error)
	LogSetLevel(context.Context, string, string) error
}
type APIAlg jwt.HMACSHA

var _ ICommon = &Common{}

type Common struct {
	fx.In
	APISecret *APIAlg
}

type jwtPayload struct {
	Allow []string
}

func (a *Common) AuthVerify(ctx context.Context, token string) ([]permission.Permission, error) {
	var payload jwtPayload
	if _, err := jwt.Verify([]byte(token), (*jwt.HMACSHA)(a.APISecret), &payload); err != nil {
		return nil, xerrors.Errorf("JWT Verification failed: %w", err)
	}
	return payload.Allow, nil
}

func (a *Common) AuthNew(ctx context.Context, perms []permission.Permission) ([]byte, error) {
	p := jwtPayload{
		Allow: perms, // TODO: consider checking validity
	}
	return jwt.Sign(&p, (*jwt.HMACSHA)(a.APISecret))
}

func (a *Common) Version(context.Context) (Version, error) {
	return Version{
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
