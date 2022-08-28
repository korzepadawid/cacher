package cacher

import (
	"sync"
)

// shardItem
type shardItem struct {
	content  []byte
	duration int64
}

// shard
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
