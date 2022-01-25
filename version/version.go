package version

import (
	"github.com/filecoin-project/venus/venus-shared/api"
)

var CurrentCommit string

// BuildVersion is the local build version, set by build system
const BuildVersion = "1.4.0-rc1"

var UserVersion = BuildVersion + CurrentCommit

// APIVersion is a semver version of the rpc api exposed
var APIVersion = api.NewVer(1, 1, 0)
