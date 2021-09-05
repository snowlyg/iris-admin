package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/multi"
)

func Init() error{
	universalOptions := &redis.UniversalOptions{
		Addrs:       strings.Split(g.CONFIG.Redis.Addr, ","),
		Password:    g.CONFIG.Redis.Password,
		PoolSize:    g.CONFIG.Redis.PoolSize,
		IdleTimeout: 300 * time.Second,
	}
	err := multi.InitDriver(
		&multi.Config{
			DriverType:      g.CONFIG.System.CacheType,
			UniversalClient: g.CACHE},
	)
	if err !=nil{
		return err
	}
	if multi.AuthDriver == nil {
		return errors.New("初始化认证驱动失败")
	}
	g.CACHE = redis.NewUniversalClient(universalOptions)
	return nil
}

// SetCache 缓存数据
func SetCache(key string, value interface{}, expiration time.Duration) error {
	err := g.CACHE.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// DeleteCache 删除缓存数据
func DeleteCache(key string) (int64, error) {
	return g.CACHE.Del(context.Background(), key).Result()
}

// GetCacheString 获取字符串类型数据
func GetCacheString(key string) (string, error) {
	return g.CACHE.Get(context.Background(), key).Result()
}

// GetCacheBytes 获取bytes类型数据
func GetCacheBytes(key string) ([]byte, error) {
	return g.CACHE.Get(context.Background(), key).Bytes()
}

// GetCacheUint 获取uint类型数据
func GetCacheUint(key string) (uint64, error) {
	return g.CACHE.Get(context.Background(), key).Uint64()
}
