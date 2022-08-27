package cacher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDjb2Hasher(t *testing.T) {
	h := newDjb2Hasher()

	t.Run("should return the same value twice", func(t *testing.T) {
		// given
		s := "Hello"
		// when
		got1 := h.sumUint64(s)
		got2 := h.sumUint64(s)
		// then
		assert.Equal(t, got1, got2)
	})

	t.Run("should return an expected value", func(t *testing.T) {
		// given
		expected := uint64(210676686969)
		s := "Hello"
		// when
		got := h.sumUint64(s)
		// then
		assert.Equal(t, got, expected)
	})

	t.Run("should return different values when strings differ by 1 character", func(t *testing.T) {
		// given
		expected1 := uint64(210676686969)
		expected2 := uint64(6952330670010)
		s1 := "Hello"
		s2 := "Hello!"
		// when
		got1 := h.sumUint64(s1)
		got2 := h.sumUint64(s2)
		// then
		assert.NotEqual(t, got1, got2)
		assert.Equal(t, expected1, got1)
		assert.Equal(t, expected2, got2)
	})

	t.Run("should be case sensitive", func(t *testing.T) {
		// given
		s1 := "Hello"
		s2 := "hello"
		// when
		got1 := h.sumUint64(s1)
		got2 := h.sumUint64(s2)
		// then
		assert.NotEqual(t, got1, got2)
	})

	t.Run("should hash empty value (with many spaces)", func(t *testing.T) {
		// given
		s := "  "
		// when
		got := h.sumUint64(s)
		// then
		assert.NotEmpty(t, got)
	})

	t.Run("should return 5381 when no chars", func(t *testing.T) {
		// given
		s := ""
		// when
		got := h.sumUint64(s)
		// then
		assert.Equal(t, uint64(5381), got)
	})
}
