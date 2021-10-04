package cache

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/snowlyg/iris-admin/server/viper_server"
)

var (
	once        sync.Once
	cacheClient redis.UniversalClient
)

// Instance 初始化缓存服务
func Instance() redis.UniversalClient {
	once.Do(func() {
		viper_server.Init(getViperConfig())
		universalOptions := &redis.UniversalOptions{
			Addrs:       strings.Split(CONFIG.Addr, ","),
			Password:    CONFIG.Password,
			PoolSize:    CONFIG.PoolSize,
			IdleTimeout: 300 * time.Second,
		}
		cacheClient = redis.NewUniversalClient(universalOptions)
	})

	return cacheClient
}

// SetCache 缓存数据
func SetCache(key string, value interface{}, expiration time.Duration) error {
	err := Instance().Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// DeleteCache 删除缓存数据
func DeleteCache(key string) (int64, error) {
	return Instance().Del(context.Background(), key).Result()
}

// GetCacheString 获取字符串类型数据
func GetCacheString(key string) (string, error) {
	value, err := GetCacheBytes(key)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// GetCacheBytes 获取bytes类型数据
func GetCacheBytes(key string) ([]byte, error) {
	return Instance().Get(context.Background(), key).Bytes()
}

// GetCacheUint 获取uint类型数据
func GetCacheUint(key string) (uint64, error) {
	return Instance().Get(context.Background(), key).Uint64()
}
