package newcache

import "context"

type CacheFlushDaemon interface {
	Done(ctx context.Context) (done <-chan interface{})
	LoadKeys(ctx context.Context, value interface{}) (cacheKeys []string)
	LoadValues(ctx context.Context) (values []interface{})
}
