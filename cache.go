package cacher

import (
	"log"
	"sync"
	"time"
)

type Cache struct {
	config *Config
	shards []*shard
	hash   hasher
}

// New initializes thread-safe, sharded,
// efficient, in-memory key-value store (Cache).
// It validates a given config, before initialization.
// Returns an error, if any.
func New(config *Config) (*Cache, error) {
	config.setDefaults()
	if err := config.valid(); err != nil {
		return nil, err
	}
	c := Cache{
		shards: initShards(config.NumberOfShards),
		config: config,
		hash:   newDjb2Hasher(),
	}
	if config.CleanupInterval != NoCleanup {
		c.runCleaner()
	}
	return &c, nil
}

// initShards initializes a slice that contains n shards.
func initShards(n int) []*shard {
	s := make([]*shard, n)
	for i := 0; i < n; i++ {
		s[i] = newShard()
	}
	return s
}

// getShardIdx calculates an index of the shard
func (c *Cache) getShardIdx(sum uint64) int {
	return int(sum % uint64(len(c.shards)))
}

// getShard returns a *shard for a given sum(hash)
func (c *Cache) getShard(sum uint64) *shard {
	idx := c.getShardIdx(sum)
	return c.shards[idx]
}

// Put puts an item into a Cache with default expiration time.
// If you want to override expiration time, look at PutWithExpiration
func (c *Cache) Put(key string, value interface{}) {
	c.PutWithExpiration(key, value, c.config.DefaultExpiration)
}

// PutWithExpiration works like Put, but you can easily override expiration time.
func (c *Cache) PutWithExpiration(key string, value interface{}, expiration time.Duration) {
	item := shardItem{
		value:      value,
		expiration: time.Now().Add(expiration).Unix(),
	}
	if expiration == NoExpiration {
		item.expiration = noExpInt64
	}
	hash := c.hash.sumUint64(key)
	sh := c.getShard(hash)
	sh.put(hash, &item)
}

// Get gets an item with the given key, otherwise returns ErrItemNotFound.
func (c *Cache) Get(key string) (interface{}, error) {
	hash := c.hash.sumUint64(key)
	sh := c.getShard(hash)
	return sh.get(hash)
}

// Delete deletes an item from the k-v store.
// If there's no such item, its no-op.
func (c *Cache) Delete(key string) {
	hash := c.hash.sumUint64(key)
	sh := c.getShard(hash)
	sh.delete(hash)
}

// Flush cleans all shards immediately.
func (c *Cache) Flush() {
	var wg sync.WaitGroup
	for _, sh := range c.shards {
		wg.Add(1)
		sh := sh
		go func() {
			defer wg.Done()
			sh.flush()
		}()
	}
	wg.Wait()
}

// deleteExpired removes all expired items
// from the Cache.
func (c *Cache) deleteExpired() {
	var wg sync.WaitGroup
	for _, sh := range c.shards {
		wg.Add(1)
		sh := sh
		go func() {
			defer wg.Done()
			sh.removeAllExpired()
		}()
	}
	wg.Wait()
}

// runCleaner triggers deleteExpired every as configured.
func (c *Cache) runCleaner() {
	ticker := time.NewTicker(c.config.CleanupInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.deleteExpired()
			}
			log.Println("Finished removing expired elements from shards")
		}
	}()
}
