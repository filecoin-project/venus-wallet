package api

import "github.com/filecoin-project/venus-wallet/api/permission"

var _ IFullAPI = &ServiceAuth{}

// full service API permissions constraints
type ServiceAuth struct {
	CommonAuth
	WalletAuth
	WalletLockAuth
	StrategyAuth
}

func PermissionedFullAPI(a IFullAPI) IFullAPI {
	var out ServiceAuth
	permission.PermissionedAny(a, &out.WalletAuth.Internal)
	permission.PermissionedAny(a, &out.CommonAuth.Internal)
	permission.PermissionedAny(a, &out.WalletLockAuth.Internal)
	permission.PermissionedAny(a, &out.StrategyAuth.Internal)
	return &out
}
