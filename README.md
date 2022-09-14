[![Go Reference](https://pkg.go.dev/badge/github.com/korzepadawid/cacher.svg)](https://pkg.go.dev/github.com/korzepadawid/cacher)
[![tests](https://github.com/korzepadawid/cacher/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/korzepadawid/cacher/actions/workflows/tests.yml)

# cacher
A sharded, concurrent key-value store (cache) library for Go, suitable for single-machine applications.

## Installation
```shell
$ go get github.com/korzepadawid/cacher 
```

## Usage
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/korzepadawid/cacher"
)

type vertex struct {
	x, y float64
}

func main() {
	c, err := cacher.New(&cacher.Config{
		// You can disable this feature, just use cacher.NoExpiration instead.
		// It's possible to override this setting whenever you want with PutWithExpiration().
		DefaultExpiration: time.Hour,
		// I used the concept of “sharding” and split a large hash table
		// into multiple partitions to localise the effects of read/write locks.
		// NumberOfShards must be greater or equal than two.
		NumberOfShards: 20,
		// It's possible to disable this feature, use cacher.NoCleanup.
		CleanupInterval: time.Second * 5,
	})

	if err != nil {
		log.Fatal(err)
	}

	// It's not mandatory, but it's recommended.
	// Store pointers, and performance will be improved.
	c.Put("p1", &vertex{x: 0.1, y: 0.2})

	item, err := c.Get("p1")

	if err != nil {
		log.Fatal(err)
	}

	v := item.(*vertex)
	fmt.Println(v) // &{0.1 0.2}
}
```
## How does it work?

### Hashing
The library uses the [dbj2 algorithm](http://www.cse.yorku.ca/~oz/hash.html) for generating string hashes.

### Sharding

Cacher uses the concept of "sharding" and splits a large hash table into multiple partitions to localise the effects of 
"read/write" locks. With many read and write requests, there's no need to block the whole data structure unnecessarily.

### Choosing a shard
`shard_id = item_hash % total_number_of_shards`
