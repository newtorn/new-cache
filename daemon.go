package newcache

import "context"

// CacheFlushDaemon represents a flush daemon will be registered,
// should do implement of all methods.
type CacheFlushDaemon interface {
	Done(ctx context.Context) (done <-chan interface{})
	LoadKeys(ctx context.Context, value interface{}) (cacheKeys []string)
	LoadValues(ctx context.Context) (values []interface{})
}
