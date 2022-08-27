package cacher

// djb2Hasher uses one of the most popular hash "function", called "dbj2".
// http://www.cse.yorku.ca/~oz/hash.html
type djb2Hasher struct{}

func newDjb2Hasher() *djb2Hasher {
	return &djb2Hasher{}
}

func (h *djb2Hasher) sumUint64(s string) uint64 {
	var hash uint64 = 5381
	for i := 0; i < len(s); i++ {
		hash = ((hash << 5) + hash) + uint64(s[i])
	}
	return hash
}
