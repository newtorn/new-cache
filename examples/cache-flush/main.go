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

	// Call Singleton for the first time will create get-started.
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
