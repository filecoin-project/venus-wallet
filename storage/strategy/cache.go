package strategy

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"sync"
)

type StrategyCache interface {
	// when then upstream task transaction remove group succeed and next func failed,
	// remove all cache to prevent data correctly
	refresh()

	set(token string, kb *storage.KeyBind)
	remove(token, address string)
	removeToken(token string)
	removeTokens(tokens []string)
	get(token, address string) *storage.KeyBind

	setBlank(token, address string)
	through(token, address string) bool
	removeBlank(token, address string)
	removeKeyBind(kb *storage.KeyBind)
	removeAddress(address string)
}

type tokenKey = string
type addressKey = string
type keyBindKey = string

func genKeyBindKey(kb *storage.KeyBind) keyBindKey {
	return kb.Address + "|" + kb.Name
}

type strategyCache struct {
	sync.RWMutex
	blank map[string]struct{}
	cache map[tokenKey]map[addressKey]*storage.KeyBind
	// for kb remove and remove all token
	kbCache   map[keyBindKey][]tokenKey
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
func (c *strategyCache) refresh() {
	c.Lock()
	defer c.Unlock()
	c.blank = make(map[string]struct{})
	c.cache = make(map[tokenKey]map[addressKey]*storage.KeyBind)
}

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

func (c *strategyCache) removeToken(token string) {
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

func (c *strategyCache) get(token, address string) *storage.KeyBind {
	c.RLock()
	defer c.RUnlock()
	mp, exist := c.cache[token]
	if !exist {
		return nil
	}
	return mp[address]
}
func (c *strategyCache) setBlank(token, address string) {
	c.Lock()
	defer c.Unlock()
	c.blank[token+address] = struct{}{}
}
func (c *strategyCache) removeBlank(token, address string) {
	c.Lock()
	defer c.Unlock()
	delete(c.blank, token+address)
}

func (c *strategyCache) through(token, address string) bool {
	c.RLock()
	defer c.RUnlock()
	_, exist := c.blank[token+address]
	return exist
}

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

func (c *strategyCache) rmKBOnly(kb *storage.KeyBind) {
	key := genKeyBindKey(kb)
	delete(c.kbCache, key)
}

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
