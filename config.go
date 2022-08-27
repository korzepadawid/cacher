package cacher

import "time"

const (
	ConfigDefaultItemSize        int64 = 1 << 10 // 1KB
	ConfigNoExpiration           int64 = -1
	ConfigDefaultNumberOfShards  int   = 10
	ConfigDefaultCleanupInterval       = time.Minute * 2
)

// Config customizes cache.
type Config struct {
	// DefaultExpiration sets default expiration for an object in your cache,
	// by default cached objects don't expire.
	// You can easily override this setting for a specific item.
	// Look at todo: add method name...
	DefaultExpiration int64

	// DefaultItemSize is a size of a cached item, given in bytes,
	// by default maximum size of an object is 1KB.
	// Must be positive.
	DefaultItemSize int64

	// NumberOfShards n shards are going to divide your hashmap into n smaller hashmaps,
	// it is helpful to reduce a number of locks on a single data structure.
	// Default value is 10, it must be greater than 1.
	NumberOfShards int

	// CleanupInterval is an interval between cache cleanups,
	// it will remove expired items from memory,
	// default cleanup interval is set to 2 minutes.
	CleanupInterval time.Duration
}
