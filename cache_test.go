package cacher

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	t.Run("should set default values and create a cache when empty config", func(t *testing.T) {
		// given
		cfg := Config{}
		// when
		c, err := New(&cfg)
		// then
		assert.NoError(t, err)
		assert.Equal(t, configDefaultExpiration, c.config.DefaultExpiration)
		assert.Equal(t, configDefaultNumberOfShards, len(c.shards))
		assert.Equal(t, configDefaultCleanupInterval, c.config.CleanupInterval)
	})

	t.Run("should create a cache when custom values in config", func(t *testing.T) {
		// given
		cfg := Config{
			DefaultExpiration: time.Hour,
			NumberOfShards:    15,
			CleanupInterval:   time.Minute * 15,
		}
		// when
		c, err := New(&cfg)
		// then
		assert.NoError(t, err)
		assert.Equal(t, time.Hour, c.config.DefaultExpiration)
		assert.Equal(t, 15, len(c.shards))
		assert.Equal(t, time.Minute*15, c.config.CleanupInterval)
	})

	t.Run("should return an error when invalid expiration date", func(t *testing.T) {
		// given
		cfg := Config{
			DefaultExpiration: -time.Millisecond * 500,
		}
		// when
		_, err := New(&cfg)
		// then
		assert.ErrorIs(t, err, ErrInvalidDefaultExpiration)
	})

	t.Run("should return an error when invalid shards count", func(t *testing.T) {
		// given
		cfg := Config{
			NumberOfShards: 1,
		}
		// when
		_, err := New(&cfg)
		// then
		assert.ErrorIs(t, err, ErrInvalidNumberOfShards)
	})

	t.Run("should return an error when invalid cleanup interval", func(t *testing.T) {
		// given
		cfg := Config{
			CleanupInterval: -time.Minute,
		}
		// when
		_, err := New(&cfg)
		// then
		assert.ErrorIs(t, err, ErrInvalidCleanupInterval)
	})
}

func TestShardIdx(t *testing.T) {
	t.Run("should return 0 when divisible by the slice length", func(t *testing.T) {
		// given
		cfg := Config{
			NumberOfShards: 15,
		}
		c, err := New(&cfg)
		require.NoError(t, err)
		// when
		idx := c.getShardIdx(uint64(45))
		// then
		assert.Equal(t, 0, idx)
	})

	t.Run("should return valid result when not divisible by the slice length", func(t *testing.T) {
		// given
		cfg := Config{
			NumberOfShards: 15,
		}
		c, err := New(&cfg)
		require.NoError(t, err)
		// when
		idx := c.getShardIdx(uint64(16))
		// then
		assert.Equal(t, 1, idx)
	})
}

func TestCachePut(t *testing.T) {
	t.Run("should put item in cache", func(t *testing.T) {
		// given
		c, err := New(&Config{
			NumberOfShards:    10,
			DefaultExpiration: NoExpiration,
		})
		require.NoError(t, err)
		// when
		c.Put("key", "jsdfgkdfhg")
		// then
		val, err := c.Get("key")
		s, ok := val.(string)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, s, "jsdfgkdfhg")
	})

	t.Run("should put item with expiration", func(t *testing.T) {
		// given
		c, err := New(&Config{
			NumberOfShards: 10,
		})
		c.config.DefaultExpiration = -time.Minute
		require.NoError(t, err)
		// when
		c.PutWithExpiration("key", "jsdfgkdfhg", time.Hour)
		// then
		val, err := c.Get("key")
		s, ok := val.(string)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, s, "jsdfgkdfhg")
	})
}
