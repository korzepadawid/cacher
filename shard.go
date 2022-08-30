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
	// ErrItemNotFound is returned when either
	// not found or an expired item.
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

// isExpired returns true if item is expired, otherwise false.
// please note: it compares seconds, that's why it's using time.Now().Unix()
func (i *shardItem) isExpired() bool {
	if i.expiration == noExpInt64 {
		return false
	}
	return time.Now().Unix() > i.expiration
}

// shard is just a one element of a sharded hashmap.
// Contains an entries hashmap and sync.RWMutex in order
// to avoid a race condition.
type shard struct {
	entries map[uint64]*shardItem
	lock    sync.RWMutex
}

// newShard initializes and returns a new shard.
func newShard() *shard {
	return &shard{
		entries: make(map[uint64]*shardItem),
		lock:    sync.RWMutex{},
	}
}

// put puts thread-safely an item into the entries hashmap.
func (sh *shard) put(hash uint64, item *shardItem) {
	sh.lock.Lock()
	sh.entries[hash] = item
	sh.lock.Unlock()
}

// get returns the value of requested item, if there's no such item,
// or it's just expired returns ErrItemNotFound.
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

// delete deletes an item from the cache,
// if there's no such item it's no-op.
func (sh *shard) delete(hash uint64) {
	sh.lock.Lock()
	delete(sh.entries, hash)
	sh.lock.Unlock()
}

// flush flushes the entire hashmap.
func (sh *shard) flush() {
	sh.lock.Lock()
	sh.entries = make(map[uint64]*shardItem)
	sh.lock.Unlock()
}

func (sh *shard) removeAllExpired() {
	sh.lock.RLock()
	for idx, item := range sh.entries {
		sh.lock.RUnlock()
		sh.lock.Lock()
		if item.isExpired() {
			delete(sh.entries, idx)
		}
		sh.lock.Unlock()
		sh.lock.RLock()
	}
	sh.lock.RUnlock()
}
