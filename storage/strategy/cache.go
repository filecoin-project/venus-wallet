package strategy

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/filecoin-project/venus-wallet/storage"
	"sync"
)

// StrategyCache memory based wallet policy cache
type StrategyCache interface {
	// when then upstream task transaction remove group succeed and next func failed,
	// remove all cache to prevent data correctly
	refresh()
	// set cache a single keyBind with token index
	set(token string, kb *storage.KeyBind)
	// remove deletes the cache at the specified address
	remove(token, address string)
	// removeStToken deletes the cache by strategy token
	removeStToken(token string)
	// removeTokens deletes the caches that the key contain strategy tokens
	removeTokens(tokens []string)
	// get gets a keyBind by strategy token and wallet address
	get(token, address string) *storage.KeyBind
	// setBlank cache penetration data that does not exist
	setBlank(token, address string)
	// through check the token exists
	through(token, address string) bool
	// removeBlank delete  penetration data
	removeBlank(token, address string)
	// removeKeyBind delete keyBind from the cache
	removeKeyBind(kb *storage.KeyBind)
	// removeAddress deletes the cache by address
	removeAddress(address string)
}

type tokenKey = string
type addressKey = string
type keyBindKey = string

func genKeyBindKey(kb *storage.KeyBind) keyBindKey {
	return kb.Address + "|" + kb.Name
}

// strategyCache memory based wallet policy cache
type strategyCache struct {
	sync.RWMutex
	blank map[string]struct{} //prevent data penetration
	cache map[tokenKey]map[addressKey]*storage.KeyBind
	// keyBind index, for remove keyBind or token
	kbCache map[keyBindKey][]tokenKey
	// wallet address index, for remove keyBind or token
	addrCache map[addressKey][]tokenKey
}

func newStrategyCache() StrategyCache {
	return &strategyCache{
		blank:     make(map[string]struct{}),
		cache:     make(map[tokenKey]map[addressKey]*storage.KeyBind),
		addrCache: make(map[addressKey][]tokenKey),
		kbCache:   make(map[keyBindKey][]tokenKey),
	}
}

// refresh clear the cache
func (c *strategyCache) refresh() {
	c.Lock()
	defer c.Unlock()
	c.blank = make(map[string]struct{})
	c.cache = make(map[tokenKey]map[addressKey]*storage.KeyBind)
}

// set cache a single keyBind with token index
func (c *strategyCache) set(token string, kb *storage.KeyBind) {
	c.Lock()
	defer c.Unlock()
	if c.cache[token] == nil {
		c.cache[token] = make(map[addressKey]*storage.KeyBind)
	}
	c.cache[token][kb.Address] = kb

	kbKey := genKeyBindKey(kb)
	if c.kbCache[kbKey] == nil {
		c.kbCache[kbKey] = make([]string, 0, 8)
	}
	c.kbCache[kbKey] = append(c.kbCache[kbKey], token)

	if c.addrCache[kb.Address] == nil {
		c.addrCache[kb.Address] = make([]string, 0, 8)
	}
	c.addrCache[kb.Address] = append(c.addrCache[kb.Address], token)
}

// remove deletes the cache at the specified address
func (c *strategyCache) remove(token, address string) {
	c.Lock()
	defer c.Unlock()
	if c.cache[token] == nil {
		return
	}
	kb, ok := c.cache[token][address]
	if !ok {
		return
	}
	delete(c.cache[token], address)
	c.rmKBWithAddr(kb)
}

// removeStToken deletes the cache by strategy token
func (c *strategyCache) removeStToken(token string) {
	c.Lock()
	defer c.Unlock()
	addrKB, ok := c.cache[token]
	if !ok {
		return
	}
	delete(c.cache, token)
	for _, v := range addrKB {
		c.rmKBWithAddr(v)
	}
}

// removeTokens deletes the caches that the key contain strategy tokens
func (c *strategyCache) removeTokens(tokens []string) {
	c.Lock()
	defer c.Unlock()
	for _, token := range tokens {
		addrKB, ok := c.cache[token]
		if !ok {
			continue
		}
		delete(c.cache, token)
		for _, v := range addrKB {
			c.rmKBWithAddr(v)
		}
	}
}

// get gets a keyBind by strategy token and wallet address
func (c *strategyCache) get(token, address string) *storage.KeyBind {
	c.RLock()
	defer c.RUnlock()
	mp, exist := c.cache[token]
	if !exist {
		return nil
	}
	return mp[address]
}

// setBlank cache penetration data that does not exist
func (c *strategyCache) setBlank(token, address string) {
	c.Lock()
	defer c.Unlock()
	c.blank[token+address] = struct{}{}
}

// removeBlank delete  penetration data
func (c *strategyCache) removeBlank(token, address string) {
	c.Lock()
	defer c.Unlock()
	delete(c.blank, token+address)
}

// through check the token exists
func (c *strategyCache) through(token, address string) bool {
	c.RLock()
	defer c.RUnlock()
	_, exist := c.blank[token+address]
	return exist
}

// rmKBWithAddr deletes keyBind from the cache
func (c *strategyCache) rmKBWithAddr(kb *storage.KeyBind) {
	key := genKeyBindKey(kb)
	tokens, ok := c.kbCache[key]
	if !ok {
		return
	}
	delete(c.kbCache, key)
	for _, v := range tokens {
		c.rmTokenInAddrCache(v, kb.Address)
	}
}

// rmKBOnly just delete keyBind in the kbCache
func (c *strategyCache) rmKBOnly(kb *storage.KeyBind) {
	key := genKeyBindKey(kb)
	delete(c.kbCache, key)
}

// removeKeyBind delete keyBind from the cache
func (c *strategyCache) removeKeyBind(kb *storage.KeyBind) {
	c.Lock()
	defer c.Unlock()
	key := genKeyBindKey(kb)
	tokens, ok := c.kbCache[key]
	if !ok {
		return
	}
	delete(c.kbCache, key)
	for _, v := range tokens {
		if c.cache[v] != nil {
			delete(c.cache[v], kb.Address)
		}
		c.rmTokenInAddrCache(v, kb.Address)
	}
}

// rmTokenInAddrCache deletes the cache by strategy token and wallet addr
func (c *strategyCache) rmTokenInAddrCache(token, addr string) {
	tokens, ok := c.addrCache[addr]
	if !ok {
		return
	}
	linq.From(tokens).Where(func(i interface{}) bool {
		return i.(string) != token
	}).ToSlice(&tokens)
	c.addrCache[addr] = tokens
}

// removeAddress deletes the cache by address
func (c *strategyCache) removeAddress(address string) {
	c.Lock()
	defer c.Unlock()
	tokens, ok := c.addrCache[address]
	if !ok {
		return
	}
	delete(c.addrCache, address)
	for _, v := range tokens {
		if c.cache[v] != nil {
			kb, ok := c.cache[v][address]
			if !ok {
				continue
			}
			delete(c.cache[v], address)
			c.rmKBOnly(kb)
		}
	}
}
