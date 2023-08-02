# venus-wallet changelog

## v1.12.0

* fix: panic when ParseObj is nil by @diwufeiwen in https://github.com/filecoin-project/venus-wallet/pull/134
* opt: use toml lib to decode config by @simlecode in https://github.com/filecoin-project/venus-wallet/pull/138
* feat: upgrade the way of generating all permissions by @diwufeiwen in https://github.com/filecoin-project/venus-wallet/pull/139
* fix: correctly generate token of api by @diwufeiwen in https://github.com/filecoin-project/venus-wallet/pull/140
* fix: permission not found by @simlecode in https://github.com/filecoin-project/venus-wallet/pull/142

## v1.11.0

* bump up version to v1.11.0

## v1.11.0-rc1

* feat: add sign recorder by @LinZexiao /保存签名记录 [[#123](https://github.com/filecoin-project/venus-wallet/pull/123)]
* feat: add docker push by @hunjixin 增加推送到镜像仓库的功能 [[#131](https://github.com/filecoin-project/venus-wallet/pull/131)]

## v1.10.1

* 修复创建 delegated 失败 [[#128](https://github.com/filecoin-project/venus-wallet/pull/128)]

## v1.10.0

* 支持 delegated 地址 [[#119](https://github.com/filecoin-project/venus-wallet/pull/119)]
* 升级 venus-gateway 和 venus-auth 版本到 v1.10.0
* 升级 go-jsonrpc 版本到 v0.1.7
