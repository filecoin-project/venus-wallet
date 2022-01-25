package strategy

import (
	"testing"

	types "github.com/filecoin-project/venus/venus-shared/types/wallet"
	"gotest.tools/assert"
)

func TestCacheFlow(t *testing.T) {
	cache := newStrategyCache()
	tk1 := "token1"
	addr1 := "address1"

	addr2 := "address2"

	kb1 := &types.KeyBind{
		Name:    "kb1",
		Address: addr1,
	}
	cache.set(tk1, kb1)
	kbTmp1 := cache.get(tk1, addr1)
	// Check the query is correct
	assert.DeepEqual(t, kb1, kbTmp1)
	// cache tb
	// TOKEN	ADDR	KB
	// tk1   	addr1  	kb1

	tk2 := "token2"
	kb2 := &types.KeyBind{
		Name:    "kb2",
		Address: addr1,
	}
	cache.set(tk2, kb2)
	kbTmp2 := cache.get(tk2, addr1)
	// Check the query is correct
	assert.DeepEqual(t, kb2, kbTmp2)
	// cache tb
	// TOKEN	ADDR	KB
	// tk1   	addr1  	kb1
	// tk2		addr1	kb2

	tk3 := "token3"
	kb3 := &types.KeyBind{
		Name:    "kb3",
		Address: addr2,
	}
	cache.set(tk3, kb3)

	// Check the query is correct
	kbTmp3Error := cache.get(tk3, addr1)
	if kbTmp3Error != nil {
		t.Fatal("data match error")
	}
	// Check the query is correct
	kbTmp3 := cache.get(tk3, addr2)
	assert.DeepEqual(t, kb3, kbTmp3)

	// cache tb
	// TOKEN	ADDR	KB
	// tk1   	addr1  	kb1
	// tk2		addr1	kb2
	// tk3 		addr2	kb3

	kb5 := &types.KeyBind{
		Name:    "kb5",
		Address: addr2,
	}
	cache.set(tk1, kb5)
	kb5Tmp := cache.get(tk1, addr2)
	assert.DeepEqual(t, kb5, kb5Tmp)
	assert.Equal(t, len(cache.(*strategyCache).kbCache), 4)
	// cache tb
	// TOKEN	ADDR	KB
	// tk1   	addr1  	kb1
	// tk2		addr1	kb2
	// tk3 		addr2	kb3
	// tk1		addr2	kb5

	cache.remove(tk2, addr1)
	kbTmp2 = cache.get(tk2, addr1)
	if kbTmp2 != nil {
		t.Fatal("remove failed")
	}
	// cache tb
	// TOKEN	ADDR	KB
	// tk1   	addr1  	kb1
	// tk2		addr1	kb2  x
	// tk3 		addr2	kb3
	// tk1		addr2	kb5
	assert.Equal(t, len(cache.(*strategyCache).kbCache), 3)

	cache.removeAddress(addr2)

	kbTmp3 = cache.get(tk3, addr2)
	kb5Tmp = cache.get(tk1, addr2)
	assert.DeepEqual(t, kbTmp3, kb5Tmp)
	assert.Equal(t, len(cache.(*strategyCache).kbCache), 1)
	// cache tb
	// TOKEN	ADDR	KB
	// tk1   	addr1  	kb1
	// tk2		addr1	kb2  x
	// tk3 		addr2	kb3  x
	// tk1		addr2	kb5  x

	kb6 := &types.KeyBind{
		Name:    "kb6",
		Address: addr2,
	}
	cache.set(tk1, kb6)

	// cache tb
	// TOKEN	ADDR	KB
	// tk1   	addr1  	kb1
	// tk2		addr1	kb2  x
	// tk3 		addr2	kb3  x
	// tk1		addr2	kb5  x
	// tk1		addr2	kb6
	cache.removeStToken(tk1)
	kb1 = cache.get(tk1, addr1)
	kb6 = cache.get(tk1, addr2)
	assert.DeepEqual(t, kb1, kb6)
	// cache tb
	// TOKEN	ADDR	KB
	// tk1   	addr1  	kb1  x
	// tk2		addr1	kb2  x
	// tk3 		addr2	kb3  x
	// tk1		addr2	kb5  x
	// tk1		addr2	kb6  x
	assert.Equal(t, len(cache.(*strategyCache).kbCache), 0)

	KeyBinds(cache, t)
}

func KeyBinds(cache StrategyCache, t *testing.T) {
	addr := "addra"
	kbA := &types.KeyBind{
		Name:    "kb5",
		Address: addr,
	}
	cache.set("tka-1", kbA)
	cache.set("tka-2", kbA)
	cache.set("tka-3", kbA)
	cache.set("tka-4", kbA)
	cache.set("tka-5", kbA)

	cache.removeKeyBind(kbA)
	assert.Equal(t, len(cache.(*strategyCache).kbCache), 0)
}
