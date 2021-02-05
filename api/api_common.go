package api

import (
	"context"
	"fmt"
	"github.com/ipfs-force-community/venus-wallet/build"
)

type Permission = string

type Common interface {
	// Auth
	AuthVerify(ctx context.Context, token string) ([]Permission, error)
	AuthNew(ctx context.Context, perms []Permission) ([]byte, error)

	// Version provides information about API provider
	Version(context.Context) (Version, error)

	LogList(context.Context) ([]string, error)
	LogSetLevel(context.Context, string, string) error
}

// Version provides various build-time information
type Version struct {
	Version string

	// APIVersion is a binary encoded semver version of the remote implementing
	// this api
	//
	// See APIVersion in build/version.go
	APIVersion build.Version
}

func (v Version) String() string {
	return fmt.Sprintf("%s+api%s", v.Version, v.APIVersion.String())
}
