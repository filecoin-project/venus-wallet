package sqlite

import (
	"testing"
	"time"

	assert2 "gotest.tools/assert"

	"golang.org/x/exp/rand"

	"github.com/filecoin-project/venus/venus-shared/types"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/venus-wallet/storage"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/crypto/aes"
)

func setup(t *testing.T) storage.KeyStore {
	conn, err := NewDB(&config.DBConfig{
		Conn: "file::memory:",
	})
	assert.NoError(t, err)
	return NewKeyStore(conn)
}

func randBytes(t *testing.T, length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	assert.NoError(t, err)
	return bytes
}

func mockData(t *testing.T, keyStore storage.KeyStore, addr string) ([]byte, []byte) {
	password := randBytes(t, 10)
	data := randBytes(t, 5)
	crypto, err := aes.EncryptData(password, data, 2, 2)
	assert.NoError(t, err)
	key := &aes.EncryptedKey{
		Address: addr,
		KeyType: types.KTBLS,
		Crypto:  crypto,
	}
	err = keyStore.Put(key)
	assert.NoError(t, err)
	return password, data
}

func Test_sqliteStorage_PutAndList(t *testing.T) {
	keyStore := setup(t)
	rand.Seed(uint64(time.Now().Unix()))

	addr := "f3uyk4vweulsdbeqfnx7g4swk2zaa4p5xnmcuqvecyuwoggvlfagruxippti2v7sc2lzyop72pyrkr2ks2xc7q"
	mockData(t, keyStore, addr)
	mockData(t, keyStore, addr)
	addr2 := "f12b5jp4z7zqdiogs7n2hpqgknxiazubl426il5xi"
	mockData(t, keyStore, addr2)

	msgs, err := keyStore.List()
	assert.NoError(t, err)
	assert.Len(t, msgs, 2)
}

func Test_sqliteStorage_HasGet(t *testing.T) {
	keyStore := setup(t)
	rand.Seed(uint64(time.Now().Unix()))

	addr, _ := address.NewFromString("f3uyk4vweulsdbeqfnx7g4swk2zaa4p5xnmcuqvecyuwoggvlfagruxippti2v7sc2lzyop72pyrkr2ks2xc7q")
	k1pass, k1data := mockData(t, keyStore, addr.String())

	addr2, _ := address.NewFromString("f12b5jp4z7zqdiogs7n2hpqgknxiazubl426il5xi")
	k2pass, k2data := mockData(t, keyStore, addr2.String())

	key1Get, err := keyStore.Get(addr)
	assert.NoError(t, err)
	data, err := aes.Decrypt(key1Get.Crypto, k1pass)
	assert.NoError(t, err)
	assert2.DeepEqual(t, k1data, data)

	has, err := keyStore.Has(addr)
	assert.True(t, has)
	assert.NoError(t, err)

	key2Get, err := keyStore.Get(addr2)
	assert.NoError(t, err)
	data2, err := aes.Decrypt(key2Get.Crypto, k2pass)
	assert.NoError(t, err)
	assert2.DeepEqual(t, k2data, data2)

	has, err = keyStore.Has(addr2)
	assert.True(t, has)
	assert.NoError(t, err)

	addrNotFound, _ := address.NewFromString("f3vno7td7s767d55yij3lucf5z3jvk2x7mwgzlbw7mdlcfxlwtzozcit6kswmmfmlcc7evtopthnkb32q6n2xa")
	_, err = keyStore.Get(addrNotFound)
	assert.EqualError(t, err, "record not found")

	has, err = keyStore.Has(addrNotFound)
	assert.False(t, has)
	assert.NoError(t, err)
}

func Test_sqliteStorage_Delete(t *testing.T) {
	keyStore := setup(t)
	rand.Seed(uint64(time.Now().Unix()))

	addr, _ := address.NewFromString("f3uyk4vweulsdbeqfnx7g4swk2zaa4p5xnmcuqvecyuwoggvlfagruxippti2v7sc2lzyop72pyrkr2ks2xc7q")
	mockData(t, keyStore, addr.String())

	addr2, _ := address.NewFromString("f12b5jp4z7zqdiogs7n2hpqgknxiazubl426il5xi")
	mockData(t, keyStore, addr2.String())

	assert.NoError(t, keyStore.Delete(addr))

	has, err := keyStore.Has(addr)
	assert.False(t, has)
	assert.NoError(t, err)

	//confirm not delete other items
	has, err = keyStore.Has(addr2)
	assert.True(t, has)
	assert.NoError(t, err)
}
