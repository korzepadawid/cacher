# cacher

## Installation
```shell
$ go get github.com/korzepadawid/cacher 
```

## Usage
```go
import "github.com/korzepadawid/cacher"

type vertex struct {
    x, y float64
}

c, err := cacher.New(&cacher.Config{
    DefaultExpiration: time.Hour,
    NumberOfShards:    20,
    CleanupInterval:   time.Second * 5,
})

if err != nil {
    // handle me
}

c.Put("p1", &vertex{x: 0.1, y: 0.2})
v, err := c.Get("p1")

if err != nil {
    // handle me
}

p := v.(*vertex)
fmt.Println(p) // &{0.1 0.2}
```