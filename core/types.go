package core

import (
	"crypto/rand"
	"io"
	"io/ioutil"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

const StringEmpty = ""

var RandSignBytes = func() []byte {
	randSignBytes, err := ioutil.ReadAll(io.LimitReader(rand.Reader, 32))
	if err != nil {
		panic(xerrors.Errorf("rand secret failed %v", err))
	}

	return randSignBytes
}()

type Address = address.Address
type Signature = crypto.Signature
type MethodNum = abi.MethodNum
type Cid = cid.Cid

var (
	NilAddress = Address{}
)
