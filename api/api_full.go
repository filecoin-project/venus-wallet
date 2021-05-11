package api

import "github.com/filecoin-project/venus-wallet/api/permission"

var _ IFullAPI = &ServiceAuth{}

// full service API permissions constraints
type ServiceAuth struct {
	CommonAPIAdapter
	WalletAPIAdapter
	WalletLockAPIAdapter
	StrategyAPIAdapter
}

func PermissionedFullAPI(a IFullAPI) IFullAPI {
	var out ServiceAuth
	permission.PermissionedAny(a, &out.WalletAPIAdapter.Internal)
	permission.PermissionedAny(a, &out.CommonAPIAdapter.Internal)
	permission.PermissionedAny(a, &out.WalletLockAPIAdapter.Internal)
	permission.PermissionedAny(a, &out.StrategyAPIAdapter.Internal)
	return &out
}
