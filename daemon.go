package newcache

import "context"

// CacheFlushDaemon represents a flush daemon will be registered,
// should do implement of all methods.
type CacheFlushDaemon interface {
	// Done tells the cache service quit this cache flush daemon.
	Done(ctx context.Context) (done <-chan interface{})
	// LoadKeys loads keys with current value by cache daemon setting.
	LoadKeys(ctx context.Context, value interface{}) (cacheKeys []string)
	// LoadValues loads values by cache daemon setting.
	LoadValues(ctx context.Context) (values []interface{})
}
