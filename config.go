package cacher

import (
	"errors"
	"time"
)

const (
	NoExpiration = time.Duration(-1)
	NoCleanup    = time.Duration(-1)

	configDefaultExpiration      = NoExpiration
	configDefaultNumberOfShards  = 10
	configDefaultCleanupInterval = time.Minute * 2
)

var (
	ErrInvalidDefaultExpiration = errors.New("invalid default expiration")
	ErrInvalidNumberOfShards    = errors.New("number of shards must be greater than 1")
	ErrInvalidCleanupInterval   = errors.New("cleanup interval must be greater than 0")
)

// Config customizes cache.
type Config struct {
	// DefaultExpiration sets default expiration time for an item in a cache,
	// you can easily override this setting whenever you put new item to the cache.
	// Value must be greater than zero.
	// Default value is NoExpiration.
	DefaultExpiration time.Duration

	// NumberOfShards n shards are going to divide your hashmap into n smaller hashmaps,
	// it is helpful to reduce a number of locks on a single data structure.
	// Default value is 10, it must be greater than 1.
	NumberOfShards int

	// CleanupInterval is an interval between cache cleanups,
	// it will remove expired items from memory,
	// default cleanup interval is set to 2 minutes.
	// You can disable cleanups with NoCleanup
	CleanupInterval time.Duration
}

// valid validates existing config, returns an error
func (c *Config) valid() error {
	if c.DefaultExpiration != NoExpiration && c.DefaultExpiration < 1 {
		return ErrInvalidDefaultExpiration
	}
	if c.NumberOfShards < 2 {
		return ErrInvalidNumberOfShards
	}
	if c.CleanupInterval != NoCleanup && c.CleanupInterval < 1 {
		return ErrInvalidCleanupInterval
	}
	return nil
}

// setDefaults sets default values on empty fields
func (c *Config) setDefaults() {
	if c.DefaultExpiration == 0 {
		c.DefaultExpiration = configDefaultExpiration
	}
	if c.NumberOfShards == 0 {
		c.NumberOfShards = configDefaultNumberOfShards
	}
	if c.CleanupInterval == 0 {
		c.CleanupInterval = configDefaultCleanupInterval
	}
}
