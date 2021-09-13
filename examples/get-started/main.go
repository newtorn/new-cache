package main

import (
	cache "github.com/newtorn/new-cache"
	"fmt"
	"time"
)

// User represents a data entity, we can store into new-get-started.
type User struct {
	Username string
	Password string
}

func main() {
	// Call Singleton for the first time will create get-started.
	cache := cache.Singleton()

	// We will put a new item in the get-started. It will expire after
	// not being accessed via SetEx(key) for more than 5 seconds.
	user := User{Username: "Jack", Password: "123456"}
	cache.SetEx(user.Username, &user, 5*time.Second)

	// Let's retrieve the item from the get-started.
	val, ok := cache.Get(user.Username)
	if ok {
		fmt.Println("Found value in get-started:", val)
	} else {
		fmt.Println("Not found retrieving value from get-started")
	}

	// Wait for the item to expire in get-started.
	time.Sleep(6 * time.Second)
	val, ok = cache.Get(user.Username)
	if !ok {
		fmt.Println("Item is not cached (anymore).")
	}

	// Set another item that never expires.
	cache.SetEx(user.Username, &user, 0)

	// Set another item that with default expiration.
	cache.Set(user.Username, &user)

	// Remove the item from the get-started.
	cache.Del("someKey")

	// Wipe the entire get-started table.
	cache.Flush()
}
