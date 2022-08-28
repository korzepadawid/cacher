package cacher

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
func initShards(n uint32) []*shard {
	s := make([]*shard, n)
	for i := uint32(0); i < n; i++ {
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
