package db_proc

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/migrations"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/filecoin-project/go-state-types/crypto"
	api2 "github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/lib/sigs"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/v2/actors/builtin/market"
	"github.com/filecoin-project/specs-actors/v2/actors/builtin/paych"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

type Types struct {
	Type      reflect.Type
	signBytes func(i interface{}) ([]byte, error)
}

const SAVE_OBJECT_JSON_FMT = true

var SupportedMsgTypes = map[api.MsgType]*Types{
	api.MTDealProposal: {reflect.TypeOf(market.DealProposal{}), func(i interface{}) ([]byte, error) {
		return cborutil.Dump(i)
	}},
	api2.MTDrawRandomParam: {reflect.TypeOf(api2.DrawRandomParams{}), func(in interface{}) ([]byte, error) {
		param := in.(*api2.DrawRandomParams)
		return param.SignBytes()
	}},
	api2.MTSignedVoucher: {reflect.TypeOf(paych.SignedVoucher{}), func(in interface{}) ([]byte, error) {
		return (in.(*paych.SignedVoucher)).SigningBytes()
	}},
	api2.MTStorageAsk: {reflect.TypeOf(storagemarket.StorageAsk{}), func(in interface{}) ([]byte, error) {
		return cborutil.Dump(in)
	}},
	api2.MTAskResponse: {reflect.TypeOf(network.AskResponse{}), func(in interface{}) ([]byte, error) {
		newAsk := in.(*network.AskResponse).Ask.Ask
		oldAsk := &migrations.StorageAsk0{newAsk.Price, newAsk.VerifiedPrice, newAsk.MinPieceSize,
			newAsk.MaxPieceSize, newAsk.Miner, newAsk.Timestamp, newAsk.Expiry, newAsk.SeqNo}
		return cborutil.Dump(oldAsk)
	}},
	api2.MTNetWorkResponse: {reflect.TypeOf(network.Response{}), func(in interface{}) ([]byte, error) {
		return cborutil.Dump(in)
	}},
	// api2.MTProviderDealState: {reflect.TypeOf(storagemarket.ProviderDealState{}), func(in interface{}) ([]byte, error) {
	// 	return cborutil.Dump(in)
	// }},
	api2.MTClientDeal: {reflect.TypeOf(market.ClientDealProposal{}), func(in interface{}) ([]byte, error) {
		ni, err := cborutil.AsIpld(in)
		if err != nil {
			return nil, err
		}
		return ni.Cid().Bytes(), nil
	}},
	api.MTBlock: {reflect.TypeOf(types.BlockHeader{}), func(in interface{}) ([]byte, error) {
		return in.(*types.BlockHeader).SigningBytes()
	}},
	api.MTChainMsg: {reflect.TypeOf(types.Message{}), func(in interface{}) ([]byte, error) {
		msg := in.(*types.Message)
		return msg.Cid().Bytes(), nil
	}},
}

func getSignBytes(toSign []byte, meta api.MsgMeta) (interface{}, []byte, error) {
	t := SupportedMsgTypes[meta.Type]
	if t == nil {
		return nil, nil, fmt.Errorf("unsupported msgtype:%s", meta.Type)
	}
	in := reflect.New(t.Type).Interface()
	unmarshaler, isok := in.(cbor.Unmarshaler)
	if !isok {
		return nil, nil, fmt.Errorf("type:%s is is not an 'unmarhsaler'", t.Type.Name())
	}
	if err := unmarshaler.UnmarshalCBOR(bytes.NewReader(toSign)); err != nil {
		return nil, nil, xerrors.Errorf("cborunmarshal to %s failed:%w", t.Type.Name(), err)
	}

	var data []byte
	data, err := t.signBytes(in)

	return in, data, err
}

func (dp *DbProc) signMsg(cidBytes, meta []byte) (*crypto.Signature, error) {
	var wallet, tx, key, signature = &Wallet{}, (*gorm.DB)(nil), (*Key)(nil), (*crypto.Signature)(nil)
	var msg, err = types.DecodeMessage(meta)

	if nil != err {
		return nil, xerrors.Errorf("decode message failed:%w", err)
	}

	tx = dp.Db.Begin()

	defer func() {
		if rbkErr := tx.Rollback().Error; (rbkErr!=nil && rbkErr != sql.ErrTxDone) {
			loger.Warnf("signMsg rollback failed:%s", rbkErr.Error())
		}
	}()

	tx.Model(wallet).Where("address = ?", msg.From.String())
	if msg.Nonce == 0 {
		tx.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err = tx.First(&wallet).Error; err != nil {
		return nil, xerrors.Errorf("find wallet(%s) failed:%w", err)
	}

	var txId cid.Cid
	if msg.Nonce == wallet.Nonce || msg.Nonce == 0 {
		msg.Nonce = wallet.Nonce
		wallet.Nonce++
		if err = tx.Table("wallets").Where("address=?", wallet.Address).
			Update("nonce", wallet.Nonce).Error; err != nil {
			return nil, xerrors.Errorf("update wallet(%s) nonce(%d) failed:%w",
				wallet.Address, wallet.Nonce)
		}
		txId = msg.Cid()
	} else if msg.Nonce < wallet.Nonce { // check cid, if there is a customized 'nonce'
		txId = msg.Cid()
		if bytes.Compare(txId.Bytes(), cidBytes) != 0 {
			return nil, fmt.Errorf("un-matched cid bytes")
		}
	} else {
		return nil,
			fmt.Errorf("pushing message with 'nonce' bigger than stored 'nonce' may cause un-expected situation")
	}

	if key, err = wallet.GetKey(); err != nil {
		return nil, xerrors.Errorf("wallet(%s) getkey failed:%w", err)
	}

	if signature, err = sigs.Sign(ActSigType(key.Type), key.PrivateKey, txId.Bytes()); nil != err {
		return nil, xerrors.Errorf("wallet(%s) sign message(%s) failed:%w",
			msg.From.String(), txId.String(), err)
	}

	sqlMsg := SqlChainSignedMsg{types.SignedMessage{
		Message:   *msg,
		Signature: *signature,
	}}

	var signedMsg = SignedMsg{
		Address:   msg.From.String(),
		Cid:       sqlMsg.Cid().String(),
		SignedMsg: sqlMsg,
		Nonce:     msg.Nonce,
	}

	if err = tx.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "address"}, {Name: "nonce"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at", "cid", "signed_msg"})}).
		Model(signedMsg).
		Save(&signedMsg).Error
		err != nil {
		return nil, err
	}

	return signature, tx.Commit().Error
}

func (dp *DbProc) signCborObject(signer address.Address, toSign []byte, meta api.MsgMeta) (*crypto.Signature, error) {
	var toSignData, err = []byte(nil), error(nil)
	var in interface{}
	if in, toSignData, err = getSignBytes(toSign, meta); err != nil {
		return nil, xerrors.Errorf("get sign bytes failed:%w", err)
	}

	var signature *crypto.Signature
	if signature, err = dp.signRawBytes(signer, toSignData); err != nil {
		return nil, err
	}

	var sd = &SignedData{Signer: signer.String(), Data: toSign, Meta: SqlMsgMeta(meta)}

	if SAVE_OBJECT_JSON_FMT {
		var jsonData []byte
		if jsonData, err = json.Marshal(in); err == nil {
			sd.Data = jsonData
		} else {
			loger.Errorf("json marshal(msgtype:%s) failed:%s", meta.Type, err.Error())
		}
	}

	return signature, dp.Db.Model((*SignedData)(nil)).Save(sd).Error
}

func (dp *DbProc) signRawBytes(signer address.Address, toSign []byte) (*crypto.Signature, error) {
	var err, wallet, key = error(nil), &Wallet{}, (*Key)(nil)
	if err = dp.Db.First(wallet, "address = ?", signer.String()).Error; err != nil {
		return nil, xerrors.Errorf("find wallet(%s) failed:%w", signer.String(), err)
	}

	if key, err = wallet.GetKey(); err != nil {
		return nil, xerrors.Errorf("wallet(%s) getkey failed:%w", err)
	}
	return sigs.Sign(ActSigType(key.Type), key.PrivateKey, toSign)
}
