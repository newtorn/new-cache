# NewCache

[![Latest Release](https://img.shields.io/github/release/newtorn/new-cache.svg)](https://github.com/newtorn/new-cache/releases)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/newtorn/new-cache)

Golang cache lib, easy to register your cache flush daemon.

## Installation

To install new-cache, simply run:

    go get github.com/newtorn/new-cache

To compile it from source:

    cd $GOPATH/src/github.com/newtorn/new-cache
    go get -u -v
    go build && go test -v

## Example
### Get Started
```go
package main

import (
	"fmt"
	"time"

	"github.com/newtorn/new-cache"
)

// User represents a data entity, we can store into the cache.
type User struct {
	Username string
	Password string
}

func main() {
	// Call Singleton for the first time will create cache.
	cache := newcache.Singleton()

	// We will put a new item in the cache. It will expire after
	// not being accessed via SetEx(key) for more than 5 seconds.
	user := User{Username: "Jack", Password: "123456"}
	cache.SetEx(user.Username, &user, 5*time.Second)

	// Let's retrieve the item from the cache.
	val, ok := cache.Get(user.Username)
	if ok {
		fmt.Println("Found value in cache:", val)
	} else {
		fmt.Println("Not found retrieving value from cache")
	}

	// Wait for the item to expire in cache.
	time.Sleep(6 * time.Second)
	val, ok = cache.Get(user.Username)
	if !ok {
		fmt.Println("Item is not cached (anymore).")
	}

	// Set another item that never expires.
	cache.SetEx(user.Username, &user, 0)

	// Set another item that with default expiration.
	cache.Set(user.Username, &user)

	// Remove the item from the cache.
	cache.Del("someKey")

	// Wipe the entire cache table.
	cache.Flush()
}
```

To run this example, go to examples/get-started/ and run:

    go run main.go

### Cache Register Flush Daemon
```go
package main

import (
	"context"
	"fmt"
	"time"

	cache "github.com/newtorn/new-cache"
)

// User represents a data entity, we can store into cache.
type User struct {
	Username string
	Password string
}

// userService implements a CacheFLushDaemon for cache registering.
type userService struct {
	users []*User
	done  chan interface{}
}

func (us *userService) Done(ctx context.Context) (done <-chan interface{}) {
	return us.done
}

func (us *userService) LoadKeys(ctx context.Context, value interface{}) []string {
	return []string{
		value.(*User).Username,
	}
}

func (us *userService) LoadValues(ctx context.Context) []interface{} {
	values := make([]interface{}, len(us.users))
	for i := 0; i < len(us.users); i++ {
		values[i] = us.users[i]
	}
	return values
}

func main() {
	ctx := context.Background()
	// Only init once new-cache configuration, must be called before singleton and register.
	// cache.InitOnce(cache.CacheConfig{})
	cache.InitOnce(cache.CacheConfig{
		DefaultExpiration: time.Duration(5) * time.Minute,
		CleanupInterval:   time.Duration(2) * time.Minute,
		FlushTimerTime:    time.Duration(30) * time.Second,
	})

	// Call Singleton for the first time will create cache.
	cache := cache.Singleton()

	// Create a user service as a flush daemon.
	// And register it into cache daemons.
	us := userService{
		done: make(chan interface{}),
		users: []*User{
			{"Jack", "123"},
			{"Jane", "123456"},
			{"Mark", "abd123"},
			{"Michale", "123aba"},
		},
	}
	cache.Register(ctx, &us)

	user, ok := cache.Get("Jack")
	if ok {
		fmt.Println(user.(*User))
	}
}
```

To run this example, go to examples/cache-flush/ and run:

    go run main.go


You can find a [few more examples here](https://github.com/newtorn/new-cache/tree/master/examples). Also see our
test-cases in cache_test.go for further working examples.