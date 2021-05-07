package common

import (
	"fmt"
	"github.com/filecoin-project/venus-wallet/version"
)

// Version provides various build-time information
type Version struct {
	Version string

	// APIVersion is a binary encoded semver version of the remote implementing
	// this api
	//
	// See APIVersion in build/version.go
	APIVersion version.Version
}

func (v Version) String() string {
	return fmt.Sprintf("%s+api%s", v.Version, v.APIVersion.String())
}
