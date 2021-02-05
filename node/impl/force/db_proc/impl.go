package db_proc

import (
	"context"
	"database/sql"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	mtype "github.com/ipfs-force-community/venus-wallet/chain/types"
	"github.com/ipfs-force-community/venus-wallet/lib/sigs"
	"github.com/ipfs-force-community/venus-wallet/node/config"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
	"gorm.io/gorm/clause"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var loger = logging.Logger("sqlite")

type DbProcInterface interface {
	WalletPut(keyType types.KeyType) (address.Address, error)
	WalletHas(address.Address) (bool, error)
	WalletList() ([]address.Address, error)
	WalletNonceSet(address.Address, uint64) (uint64, error)
	WalletQuery(addr address.Address) (*Key, error)
	WalletExport(addr address.Address) (*types.KeyInfo, error)
	WalletImport(*types.KeyInfo) (address.Address, error)
	WalletSign(context.Context, address.Address, []byte, api.MsgMeta) (*crypto.Signature, error)
	WalletDel(address.Address) (bool, error)
	MessageResult(cid.Cid, uint64) error
	MessageDel(address.Address, uint64, uint64) (bool, error)
	MessageQuery(address.Address, uint64, uint64) ([]mtype.SignedMsg, error)
	QueryNoPubMsg() ([]SignedMsg, error)
}

type DbProc struct {
	Db *gorm.DB
}

var _ DbProcInterface = (*DbProc)(nil)

func NewDbProc(cfg *config.DbCfg) (DbProcInterface, error) {
	var db, err = gorm.Open(sqlite.Open(cfg.Conn), &gorm.Config{})
	var sqldb *sql.DB
	if err != nil {
		return nil, xerrors.Errorf("open database(%s) failed:%w", cfg.Conn, err)
	}

	if sqldb, err = db.DB(); err != nil {
		return nil, xerrors.Errorf("sqlDb failed, %w", err)
	}

	sqldb.SetConnMaxIdleTime(300)
	sqldb.SetMaxIdleConns(8)
	sqldb.SetMaxOpenConns(64)

	db = db.Debug()
	if err = db.AutoMigrate(&Wallet{}, &SignedMsg{}, &SignedData{}); err != nil {
		return nil, xerrors.Errorf("migrate failed:%w", err)
	}

	loger.Info("init db success! ...")

	return &DbProc{Db: db}, err
}

func kstoreSigType(typ crypto.SigType) types.KeyType {
	switch typ {
	case crypto.SigTypeBLS:
		return types.KTBLS
	case crypto.SigTypeSecp256k1:
		return types.KTSecp256k1
	default:
		return ""
	}
}

func ActSigType(typ types.KeyType) crypto.SigType {
	switch typ {
	case types.KTBLS:
		return crypto.SigTypeBLS
	case types.KTSecp256k1:
		return crypto.SigTypeSecp256k1
	default:
		return 0
	}
}

func NewKey(keyinfo types.KeyInfo) (*Key, error) {
	k := &Key{
		KeyInfo: keyinfo,
	}

	var err error
	k.PublicKey, err = sigs.ToPublic(ActSigType(k.Type), k.PrivateKey)
	if err != nil {
		return nil, err
	}

	switch k.Type {
	case types.KTSecp256k1:
		k.Address, err = address.NewSecp256k1Address(k.PublicKey)
		if err != nil {
			return nil, xerrors.Errorf("converting Secp256k1 to address: %w", err)
		}
	case types.KTBLS:
		k.Address, err = address.NewBLSAddress(k.PublicKey)
		if err != nil {
			return nil, xerrors.Errorf("converting BLS to address: %w", err)
		}
	default:
		return nil, xerrors.Errorf("unknown key type")
	}
	return k, nil

}

func GenerateKey(typ crypto.SigType) (*Key, error) {
	pk, err := sigs.Generate(typ)
	if err != nil {
		return nil, err
	}
	ki := types.KeyInfo{
		Type:       kstoreSigType(typ),
		PrivateKey: pk,
	}
	return NewKey(ki)
}

func (dp *DbProc) QueryNoPubMsg() ([]SignedMsg, error) {
	ms := []SignedMsg{}

	err := dp.Db.Raw("SELECT id,address,`cid`,`signed_msg`,nonce FROM `signed_msgs` WHERE `epoch`=0 ORDER BY `nonce`;").Scan(&ms).Error
	if err != nil {
		loger.Errorf("query msg which have not on chain err: %s", err)
		return nil, err
	}

	return ms, err
}

func (dp *DbProc) WalletPut(keyType types.KeyType) (address.Address, error) {
	var k, err = GenerateKey(ActSigType(keyType))
	if err != nil {
		return address.Undef, err
	}
	wallet := Wallet{
		Address: k.Address.String(),
		KeyInfo: (SqlKeyInfo)(k.KeyInfo),
		Nonce:   0,
	}
	if err = dp.Db.First(&wallet, "address=?", wallet.Address).Error; err != nil && err != gorm.ErrRecordNotFound {
		return address.Undef, err
	} else if err == gorm.ErrRecordNotFound {
		return k.Address, dp.Db.Create(&wallet).Error
	}
	return k.Address, err
}

func (dp *DbProc) WalletHas(addr address.Address) (bool, error) {
	var counts int64 = 0

	err := dp.Db.Table("wallets").Where("address=?", addr.String()).Count(&counts).Error
	if err != nil {
		return false, err
	}

	return counts > 0, err
}

func (dp *DbProc) WalletList() ([]address.Address, error) {
	loger.Infof("wallet list")
	ws := []Wallet{} // 不能是unaddressable

	err := dp.Db.Table("wallets"). /*.Select("address", "nonce").*/ Scan(&ws).Error
	if err != nil {
		return nil, err
	}

	addresses := make([]address.Address, len(ws))
	for idx, val := range ws {
		addresses[idx], _ = address.NewFromString(val.Address)
	}
	return addresses, err
}

func (dp *DbProc) WalletNonceSet(addr address.Address, nonce uint64) (uint64, error) {
	err := dp.Db.Exec("UPDATE `wallets` SET `nonce`=? WHERE `address`=?;",
		nonce, addr.String()).Error

	if err != nil {
		return 0, err
	}
	return nonce, err
}

func (dp *DbProc) WalletQuery(addr address.Address) (*Key, error) {
	res := &Wallet{}
	if err := dp.Db.Where("address=?", addr.String()).First(res).Error; err != nil {
		return nil, err
	}
	return NewKey(types.KeyInfo(res.KeyInfo))
}

func (dp *DbProc) WalletExport(addr address.Address) (*types.KeyInfo, error) {
	k, err := dp.WalletQuery(addr)
	if err != nil {
		return nil, xerrors.Errorf("failed to find key to export: %w", err)
	}

	return &k.KeyInfo, nil
}

func (dp *DbProc) WalletImport(ki *types.KeyInfo) (address.Address, error) {
	var k, err = NewKey(*ki)
	if err != nil {
		return address.Undef, err
	}

	wallet := &Wallet{
		Address: k.Address.String(),
		KeyInfo: (SqlKeyInfo)(k.KeyInfo),
	}

	if err = dp.Db.Model(&Wallet{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"private_key"}),
	}).Save(wallet).Error; err != nil {
		return address.Undef, err
	}

	return k.Address, err
}

func (dp *DbProc) findKey(addr address.Address) (*Key, error) {
	res := &Wallet{}
	if err := dp.Db.Where("address=?", addr.String()).First(res).Error; err != nil {
		return nil, err
	}
	return res.GetKey()
}

func (dp *DbProc) WalletSign(ctx context.Context, signer address.Address,
	msg []byte, meta api.MsgMeta) (*crypto.Signature, error) {
	if meta.Type == api.MTChainMsg {
		return dp.signMsg(msg, meta.Extra)
	}
	return dp.signCborObject(signer, msg, meta)
}

func (dp *DbProc) WalletDel(addr address.Address) (bool, error) {
	var err error = nil

	tmpDb := dp.Db.Table("wallets").Delete(nil, "address = ?", addr.String())

	if err = tmpDb.Error; err != nil {
		// may be it isn't explicit, but acceptable
		return false, xerrors.Errorf("delete wallet(%s) failed:%w",
			addr.String(), err)
	}

	return tmpDb.RowsAffected > 0, nil
}

func (dp *DbProc) MessageResult(cid cid.Cid, epoch uint64) error {
	var err error = nil
	err = dp.Db.Exec("UPDATE `signed_msgs` SET epoch=? WHERE cid=?;",
		epoch, cid.String()).Error

	return err
}

func (dp *DbProc) MessageDel(addr address.Address, from uint64, to uint64) (bool, error) {
	var err error = nil
	if from == to {
		err = dp.Db.Exec("DELETE FROM `signed_msgs` WHERE address=? AND nonce=?;",
			addr.String(), from).Error
	} else {
		err = dp.Db.Exec("DELETE FROM `signed_msgs` WHERE address=? AND nonce>=? AND nonce<=?;",
			addr.String(), from, to).Error
	}

	return err == nil, err
}

func (dp *DbProc) MessageQuery(addr address.Address, from uint64, to uint64) ([]mtype.SignedMsg, error) {
	var err error = nil
	ms := []SignedMsg{}
	if from == to {
		err := dp.Db.Raw("SELECT cid,nonce,signed_msg,gr_epoch,epoch FROM `signed_msgs` WHERE address=? AND nonce=?;",
			addr.String(), from).Scan(&ms).Error
		if err != nil {
			loger.Errorf("query msg err: %s", err)
			return nil, err
		}
	} else {
		err := dp.Db.Raw("SELECT cid,nonce,signed_msg,gr_epoch,epoch FROM `signed_msgs` WHERE address=? AND nonce>=? AND nonce<=?;",
			addr.String(), from, to).Scan(&ms).Error
		if err != nil {
			loger.Errorf("query msg err: %s", err)
			return nil, err
		}
	}

	res := make([]mtype.SignedMsg, len(ms))
	for i, msg := range ms {
		res[i].Cid = msg.Cid
		res[i].Nonce = msg.Nonce
		res[i].Epoch = msg.Epoch
		res[i].SignedMsg = msg.SignedMsg.Message
	}
	return res, err
}
