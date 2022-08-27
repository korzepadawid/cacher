package cacher

import (
	"errors"
	"time"
)

const (
	ConfigDefaultMaxItemSize     uint64 = 1 << 10 // 1KB
	ConfigNoExpiration           uint64 = 0
	ConfigDefaultNumberOfShards  uint32 = 10
	ConfigDefaultCleanupInterval        = time.Minute * 2
)

var (
	ErrInvalidNumberOfShards  = errors.New("number of shards must be greater than 1")
	ErrInvalidCleanupInterval = errors.New("cleanup interval must be greater than 0")
)

// Config customizes cache.
type Config struct {
	// DefaultExpiration sets default expiration for an object in your cache,
	// by default cached objects don't expire.
	// You can easily override this setting for a specific item.
	// Look at todo: add method name...
	DefaultExpiration uint64

	// DefaultMaxItemSize is a size of a cached item, given in bytes,
	// by default maximum size of an object is 1KB.
	DefaultMaxItemSize uint64

	// NumberOfShards n shards are going to divide your hashmap into n smaller hashmaps,
	// it is helpful to reduce a number of locks on a single data structure.
	// Default value is 10, it must be greater than 1.
	NumberOfShards uint32

	// CleanupInterval is an interval between cache cleanups,
	// it will remove expired items from memory,
	// default cleanup interval is set to 2 minutes.
	CleanupInterval time.Duration
}

func (c *Config) valid() error {
	if c.NumberOfShards < 2 {
		return ErrInvalidNumberOfShards
	}
	if c.CleanupInterval < 1 {
		return ErrInvalidCleanupInterval
	}
	return nil
}

// setDefaults sets default values on empty fields
func (c *Config) setDefaults() {
	// we don't need to check default expiration, since default expiration is 0
	if c.DefaultMaxItemSize == 0 {
		c.DefaultMaxItemSize = ConfigDefaultMaxItemSize
	}
	if c.NumberOfShards == 0 {
		c.NumberOfShards = ConfigDefaultNumberOfShards
	}
	if c.CleanupInterval == 0 {
		c.CleanupInterval = ConfigDefaultCleanupInterval
	}
}
