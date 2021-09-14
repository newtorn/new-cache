package newcache

import "time"

var (
	// gFlushTime 本地缓存刷新时间 (unit: s)
	gFlushTime = time.Duration(60) * time.Second
	// gExpiration 本地缓存过期时间（unit：s）
	gExpiration = time.Duration(5*60) * time.Second
	// gCleanupInterval 本地缓存定时清理时间（unit：s）
	gCleanupInterval = time.Duration(10*60) * time.Second
)

// CacheConfig sets for cache service daemon.
type CacheConfig struct {
	// local cache flush time.
	FlushTimerTime    time.Duration
	// local cache cleanup time.
	CleanupInterval   time.Duration
	// local cache expiration time.
	DefaultExpiration time.Duration
}

// InitOnce inits cache service configuration only once.
func InitOnce(conf CacheConfig) {
	gInitOnce.Do(func() {
		gFlushTime = conf.FlushTimerTime
		gExpiration = conf.DefaultExpiration
		gCleanupInterval = conf.CleanupInterval
	})
}
