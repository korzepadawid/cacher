package cacher

type cache struct {
	config *Config
	shards []*shard
	hash   hasher
}

// New initializes thread-safe, sharded,
// efficient in-memory key-value store (cache).
// It validates a given config, before initialization.
func New(config *Config) (*cache, error) {
	if err := config.valid(); err != nil {
		return nil, err
	}
	config.setDefaults()
	return &cache{
		shards: initShards(config.NumberOfShards),
		config: config,
		hash:   newDjb2Hasher(),
	}, nil
}

// initShards initializes slice that contains n shards.
func initShards(n uint32) []*shard {
	s := make([]*shard, n)
	for i := uint32(0); i < n; i++ {
		s[i] = newShard()
	}
	return s
}
