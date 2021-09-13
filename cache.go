package newcache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"log"
	"sync"
	"time"
)

var (
	gInitOnce     sync.Once
	gCacheOnce    sync.Once
	gCacheService *cacheService
	gTicker       *time.Ticker
)

type CacheService interface {
	Set(k string, v interface{})
	SetEx(k string, v interface{}, ex time.Duration)
	Get(k string) (v interface{}, ok bool)
	Del(k string)
	GetValues() []interface{}
	GetByDefault(k string, v interface{}) interface{}
	Flush()
	Register(ctx context.Context, daemon CacheFlushDaemon)
}

type cacheService struct {
	cache *cache.Cache
}

func newCacheService() *cacheService {
	return &cacheService{
		cache: cache.New(gExpiration, gCleanupInterval),
	}
}

func Singleton() CacheService {
	gCacheOnce.Do(func() {
		gTicker = time.NewTicker(gFlushTime)
		gCacheService = newCacheService()
	})
	return gCacheService
}

func (c *cacheService) Set(k string, v interface{}) {
	c.cache.SetDefault(k, v)
}

func (c *cacheService) SetEx(k string, v interface{}, ex time.Duration) {
	c.cache.Set(k, v, ex)
}

func (c *cacheService) Get(k string) (v interface{}, ok bool) {
	return c.cache.Get(k)
}

func (c *cacheService) Del(k string) {
	c.cache.Delete(k)
}

func (c *cacheService) GetByDefault(k string, v interface{}) interface{} {
	if vc, ok := c.cache.Get(k); !ok {
		return vc
	}
	return v
}

func (c *cacheService) GetKeys() []string {
	items := c.cache.Items()
	keys := make([]string, len(items))
	for key, _ := range items {
		keys = append(keys, key)
	}
	return keys
}

func (c *cacheService) GetValues() []interface{} {
	items := c.cache.Items()
	values := make([]interface{}, len(items))
	for key, _ := range items {
		values = append(values, items[key].Object)
	}
	return values
}

func (c *cacheService) Flush() {
	c.cache.Flush()
}

func (c *cacheService) Register(ctx context.Context, daemon CacheFlushDaemon) {
	c.setValues(ctx, daemon)
	done := daemon.Done(ctx)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				// timed task error
				log.Println(ctx, "Timed tasks Failed to refresh get-started，because %v", err)
			}
		}()
		for {
			select {
			case <-gTicker.C:
				c.setValues(ctx, daemon)
			case <-done:
				// stop the daemon
				return
			}
		}
	}()
}

func (c *cacheService) setValues(ctx context.Context, daemon CacheFlushDaemon) {
	values := daemon.LoadValues(ctx)
	for valueIndex, _ := range values {
		keys := daemon.LoadKeys(ctx, values[valueIndex])
		for keyIndex, _ := range keys {
			c.Set(keys[keyIndex], values[valueIndex])
		}
	}
}
