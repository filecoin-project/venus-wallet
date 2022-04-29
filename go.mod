module github.com/filecoin-project/venus-wallet

go 1.16

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	github.com/BurntSushi/toml v0.3.1
	github.com/ahmetb/go-linq/v3 v3.2.0
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-cbor-util v0.0.1
	github.com/filecoin-project/go-crypto v0.0.1
	github.com/filecoin-project/go-fil-markets v1.20.1
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec
	github.com/filecoin-project/go-state-types v0.1.3
	github.com/filecoin-project/specs-actors/v2 v2.3.6
	github.com/filecoin-project/venus v1.2.4-0.20220420072943-4d565663fa60
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gbrlsnchs/jwt/v3 v3.0.1
	github.com/google/uuid v1.3.0
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/ipfs-force-community/venus-common-utils v0.0.0-20210924063144-1d3a5b30de87
	github.com/ipfs/go-cid v0.1.0
	github.com/ipfs/go-log/v2 v2.5.0
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/multiformats/go-multiaddr v0.5.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/viper v1.7.1
	github.com/supranational/blst v0.3.4
	github.com/urfave/cli/v2 v2.3.0
	go.opencensus.io v0.23.0
	go.uber.org/dig v1.13.0 // indirect
	go.uber.org/fx v1.15.0
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
	gotest.tools v2.2.0+incompatible
)

replace github.com/filecoin-project/go-jsonrpc => github.com/ipfs-force-community/go-jsonrpc v0.1.4-0.20210731021807-68e5207079bc
