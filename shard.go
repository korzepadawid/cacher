package cacher

import (
	"errors"
	"sync"
	"time"
)

const (
	noExpInt64 int64 = -1
)

var (
	ErrItemNotFound = errors.New("item not found")
)

// shardItem contains cached item
// and expiration date
type shardItem struct {

	// value recommendation: store pointers,
	// key-value store will be more efficient
	value interface{}

	// expiration note: compares seconds
	expiration int64
}

func (i *shardItem) isExpired() bool {
	if i.expiration == noExpInt64 {
		return false
	}
	return time.Now().Unix() > i.expiration
}

type shard struct {
	entries map[uint64]*shardItem
	lock    sync.RWMutex
}

func newShard() *shard {
	return &shard{
		entries: make(map[uint64]*shardItem),
		lock:    sync.RWMutex{},
	}
}

func (sh *shard) put(hash uint64, item *shardItem) {
	sh.lock.Lock()
	sh.entries[hash] = item
	sh.lock.Unlock()
}

func (sh *shard) get(hash uint64) (interface{}, error) {
	sh.lock.RLock()
	entry, ok := sh.entries[hash]
	if !ok || entry.isExpired() {
		sh.lock.RUnlock()
		return nil, ErrItemNotFound
	}
	sh.lock.RUnlock()
	return entry.value, nil
}
