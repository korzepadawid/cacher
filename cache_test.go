package cacher

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	t.Run("should create a cache when default settings", func(t *testing.T) {
		// given
		cfg := Config{}
		// when
		c, err := New(&cfg)
		// then
		assert.NoError(t, err)
		assert.Equal(t, int(ConfigDefaultNumberOfShards), len(c.shards))
		assert.Equal(t, ConfigDefaultCleanupInterval, c.config.CleanupInterval)
		assert.Equal(t, ConfigDefaultMaxItemSize, c.config.DefaultMaxItemSize)
		assert.Equal(t, ConfigNoExpiration, c.config.DefaultExpiration)
	})

	t.Run("should create a cache when valid custom settings", func(t *testing.T) {
		// given
		cfg := Config{
			DefaultExpiration:  time.Hour,
			DefaultMaxItemSize: 1 << 12,
			NumberOfShards:     20,
			CleanupInterval:    time.Minute * 15,
		}
		// when
		c, err := New(&cfg)
		// then
		assert.NoError(t, err)
		assert.Equal(t, 20, len(c.shards))
		assert.Equal(t, time.Minute*15, c.config.CleanupInterval)
		assert.Equal(t, uint64(1<<12), c.config.DefaultMaxItemSize)
		assert.Equal(t, time.Hour, c.config.DefaultExpiration)
	})

	t.Run("should return an error when not enough shards", func(t *testing.T) {
		// given
		cfg := Config{
			NumberOfShards: 1,
		}
		// when
		_, err := New(&cfg)
		// then
		assert.ErrorIs(t, err, ErrInvalidNumberOfShards)
	})

	t.Run("should return an error when not enough shards", func(t *testing.T) {
		// given
		cfg := Config{
			CleanupInterval: -1,
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
