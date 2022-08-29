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
	t.Run("should put new item when hasn't existed yet", func(t *testing.T) {
		// given
		sh := newShard()
		// when
		sh.put(fakeHash, &shardItem{value: "test", expiration: noExpInt64})
		// then
		result, ok := sh.entries[fakeHash]
		assert.True(t, ok)
		assert.NotEmpty(t, result)
		assert.Equal(t, result.expiration, noExpInt64)
	})

	t.Run("should override when item has existed", func(t *testing.T) {
		// given
		sh := newShard()
		// when
		sh.put(fakeHash, &shardItem{value: "test", expiration: noExpInt64})
		sh.put(fakeHash, &shardItem{value: "test2", expiration: noExpInt64})
		// then
		result, ok := sh.entries[fakeHash]
		assert.True(t, ok)
		assert.NotEmpty(t, result)
		assert.Equal(t, result.expiration, noExpInt64)
		assert.True(t, ok)
		s, ok := result.value.(string)
		assert.Equal(t, "test2", s)
	})
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

func TestDelete(t *testing.T) {
	t.Run("should delete item when present", func(t *testing.T) {
		// given
		sh := newShard()
		sh.entries[fakeHash] = &shardItem{value: "test", expiration: noExpInt64}
		// when
		sh.delete(fakeHash)
		// then
		_, ok := sh.entries[fakeHash]
		assert.False(t, ok)
	})

	t.Run("should no-op when isn't present", func(t *testing.T) {
		// given
		sh := newShard()
		// when
		sh.delete(fakeHash)
		// then
		_, ok := sh.entries[fakeHash]
		assert.False(t, ok)
	})
}

func TestFlush(t *testing.T) {
	t.Run("should clear hashmap when not-empty", func(t *testing.T) {
		// given
		sh := newShard()
		for i := 0; i < 20; i++ {
			sh.entries[uint64(i)] = &shardItem{value: i, expiration: noExpInt64}
		}
		// when
		sh.flush()
		// then
		assert.Equal(t, len(sh.entries), 0)
	})
}

func TestRemoveAllExpired(t *testing.T) {
	t.Run("should remove all expired elements and omit non expiring item", func(t *testing.T) {
		// given
		sh := newShard()
		for i := 0; i < 200; i++ {
			sh.entries[uint64(i)] = &shardItem{value: i, expiration: time.Now().Add(-time.Second).Unix()}
		}
		sh.entries[200] = &shardItem{value: 200, expiration: noExpInt64}
		// when
		sh.removeAllExpired()
		// then
		it := sh.entries[200]
		val, ok := it.value.(int)
		assert.Equal(t, 1, len(sh.entries))
		assert.True(t, ok)
		assert.Equal(t, 200, val)
	})
}
