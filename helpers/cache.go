package helpers

import (
	"fmt"
	"time"
	"os"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

func InitRedisCache() *cache.Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		},
	})

	cacheClient := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return cacheClient
}
