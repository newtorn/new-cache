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

type CacheConfig struct {
	FlushTimerTime    time.Duration
	CleanupInterval   time.Duration
	DefaultExpiration time.Duration
}

func InitOnce(conf CacheConfig) {
	gInitOnce.Do(func() {
		gFlushTime = conf.FlushTimerTime
		gExpiration = conf.DefaultExpiration
		gCleanupInterval = conf.CleanupInterval
	})
}
