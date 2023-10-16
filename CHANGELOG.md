# venus-wallet changelog

## v1.14.0-rc1

* chore: update deps

## v1.13.0

* Fix/sqlite migrator in [[#155](https://github.com/filecoin-project/venus-wallet/pull/155)]

## v1.13.0-rc1

### New Feature
* feat: move sign_type to venus share in [[#143](https://github.com/filecoin-project/venus-wallet/pull/143)]
* feat: add sign record config in [[#150](https://github.com/filecoin-project/venus-wallet/pull/150)]

### Documentation And Chore
* doc: 更新一些错误的命令 / update some outdated command in [[#147](https://github.com/filecoin-project/venus-wallet/pull/147)]
* chore(CI): ignore docs in [[#149](https://github.com/filecoin-project/venus-wallet/pull/149)]
* chore: merge release v1.12 in [[#146](https://github.com/filecoin-project/venus-wallet/pull/146)]
* build(deps): bump github.com/supranational/blst from 0.3.4 to 0.3.11 in [[#152](https://github.com/filecoin-project/venus-wallet/pull/152)]
* build(deps): bump github.com/libp2p/go-libp2p from 0.27.5 to 0.27.8 in [[#151](https://github.com/filecoin-project/venus-wallet/pull/151)]

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
