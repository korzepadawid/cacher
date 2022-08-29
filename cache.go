package cacher

import (
	"sync"
	"time"
)

type cache struct {
	config *Config
	shards []*shard
	hash   hasher
}

// New initializes thread-safe, sharded,
// efficient, in-memory key-value store (cache).
// It validates a given config, before initialization.
func New(config *Config) (*cache, error) {
	config.setDefaults()
	if err := config.valid(); err != nil {
		return nil, err
	}
	return &cache{
		shards: initShards(config.NumberOfShards),
		config: config,
		hash:   newDjb2Hasher(),
	}, nil
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
func (c *cache) getShardIdx(sum uint64) int {
	return int(sum % uint64(len(c.shards)))
}

// getShard returns a *shard for a given sum(hash)
func (c *cache) getShard(sum uint64) *shard {
	idx := c.getShardIdx(sum)
	return c.shards[idx]
}

func (c *cache) Put(key string, value interface{}) {
	c.PutWithExpiration(key, value, c.config.DefaultExpiration)
}

func (c *cache) PutWithExpiration(key string, value interface{}, expiration time.Duration) {
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

func (c *cache) Get(key string) (interface{}, error) {
	hash := c.hash.sumUint64(key)
	sh := c.getShard(hash)
	return sh.get(hash)
}

func (c *cache) Delete(key string) {
	hash := c.hash.sumUint64(key)
	sh := c.getShard(hash)
	sh.delete(hash)
}

func (c *cache) Flush() {
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

func (c *cache) deleteExpired() {
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
