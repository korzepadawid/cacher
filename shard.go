package cacher

import (
	"sync"
	"time"
)

// shardItem
type shardItem struct {
	content    []byte
	expiration time.Time
}

// shard
type shard struct {
	entries map[uint64]*shardItem
	lock    sync.RWMutex
}