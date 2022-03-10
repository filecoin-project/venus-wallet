package sqlite

import (
	"database/sql/driver"
	"github.com/filecoin-project/go-address"
	"golang.org/x/xerrors"
)

type shortAddress address.Address

//lint:ignore U1000 unused
type netWork address.Network

var networkPrefix = map[address.Network]string{
	address.Testnet: address.TestnetPrefix,
	address.Mainnet: address.MainnetPrefix,
}

func shortAddressFromString(s string) (shortAddress, error) {
	addr, err := address.NewFromString(s)
	if err != nil {
		if xerrors.Is(err, address.ErrUnknownNetwork) {
			addr, err = address.NewFromString(networkPrefix[address.CurrentNetwork] + s)
		}
	}
	return shortAddress(addr), err
}

// nolint
func (n netWork) prefix() string {
	return networkPrefix[address.Network(n)]
}

func (sa *shortAddress) Scan(value interface{}) error {
	var a, ok = value.([]byte)
	if !ok {
		return xerrors.New("address should be a string")
	}
	id, err := address.NewFromBytes(append([]byte{address.CurrentNetwork}, a...))
	if err != nil {
		return err
	}
	*sa = (shortAddress)(id)
	return nil
}

func (sa shortAddress) Value() (driver.Value, error) {
	return sa.String(), nil
}

func (sa shortAddress) String() string {
	return ((address.Address)(sa)).String()[1:]
}

func (sa shortAddress) Address() address.Address {
	return address.Address(sa)
}
