module github.com/ipfs-force-community/venus-wallet

go 1.14

replace github.com/supranational/blst => github.com/supranational/blst v0.2.0

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	github.com/BurntSushi/toml v0.3.1
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-cbor-util v0.0.0-20191219014500-08c40a1e63a2
	github.com/filecoin-project/go-crypto v0.0.0-20191218222705-effae4ea9f03
	github.com/filecoin-project/go-fil-markets v1.0.10
	github.com/filecoin-project/go-jsonrpc v0.1.2
	github.com/filecoin-project/go-state-types v0.0.0-20210119062722-4adba5aaea71
	github.com/filecoin-project/lotus v1.4.1
	github.com/filecoin-project/specs-actors v0.9.13
	github.com/filecoin-project/specs-actors/v2 v2.3.3
	github.com/fxamacker/cbor v1.5.1
	github.com/gbrlsnchs/jwt/v3 v3.0.0
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-ds-badger2 v0.1.1-0.20200708190120-187fc06f714e
	github.com/ipfs/go-fs-lock v0.0.6
	github.com/ipfs/go-log v1.0.4
	github.com/ipfs/go-log/v2 v2.1.2-0.20200626104915-0016c0b4b3e4
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/multiformats/go-base32 v0.0.3
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/multiformats/go-multiaddr-dns v0.2.0
	github.com/multiformats/go-multiaddr-net v0.2.0
	github.com/stretchr/testify v1.6.1
	github.com/supranational/blst v0.1.1
	github.com/urfave/cli/v2 v2.3.0
	go.opencensus.io v0.22.6
	go.uber.org/fx v1.13.1
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.12
)
