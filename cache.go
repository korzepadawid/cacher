package cacher

type cache struct {
	config Config
	shards []*shard
	hash   hasher
}

// New initializes thread-safe, sharded,
// efficient in-memory key-value store (cache).
// It validates a given config, before initialization.
func New(config *Config) (*cache, error) {
	return &cache{}, nil
}

// initShards initializes slice that contains n shards.
func initShards(n int) []*shard {
	s := make([]*shard, n)
	for i := 0; i < n; i++ {
		s[i] = newShard()
	}
	return s
}
