package cacher

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	fakeHash = uint64(123123)
)

func TestPut(t *testing.T) {
	// given
	sh := newShard()
	// when
	sh.put(fakeHash, &shardItem{value: "test", expiration: noExpInt64})
	// then
	result, ok := sh.entries[fakeHash]
	assert.True(t, ok)
	assert.NotEmpty(t, result)
	assert.Equal(t, result.expiration, noExpInt64)
}

func TestGet(t *testing.T) {
	t.Run("should return an entry when non-expiring", func(t *testing.T) {
		// given
		sh := newShard()
		sh.entries[fakeHash] = &shardItem{value: "test", expiration: noExpInt64}
		// when
		item, err := sh.get(fakeHash)
		s, ok := item.(string)
		// then
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "test", s)
	})

	t.Run("should return an entry when expires in future", func(t *testing.T) {
		// given
		sh := newShard()
		sh.entries[fakeHash] = &shardItem{value: "test", expiration: time.Now().Add(time.Hour).Unix()}
		// when
		item, err := sh.get(fakeHash)
		s, ok := item.(string)
		// then
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "test", s)
	})

	t.Run("should return error when expired item (not found err)", func(t *testing.T) {
		// given
		sh := newShard()
		sh.entries[fakeHash] = &shardItem{value: "test", expiration: time.Now().Add(-time.Second).Unix()}
		// when
		_, err := sh.get(fakeHash)
		// then
		assert.ErrorIs(t, err, ErrItemNotFound)
	})

	t.Run("should return error when entry not found", func(t *testing.T) {
		// given
		sh := newShard()
		// when
		_, err := sh.get(fakeHash)
		// then
		assert.ErrorIs(t, err, ErrItemNotFound)
	})
}
