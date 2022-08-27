package cacher

// hasher generates unsigned int64
// value from string (generates hash key)
//
// Why uint64? generated value is always non-negative,
// there is no need to store a sign on the first bit of the result,
// int64 since we want to have fewer collisions to deal with.
type hasher interface {
	sumUint64(s string) uint64
}
