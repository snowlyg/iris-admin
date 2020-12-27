package cache

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func RedisLock(key, value string, expireTime time.Duration) bool {
	return GetRedisClusterClient().SetNX(key, value, int(expireTime.Seconds()))
}
func RedisUnLock(key, value string) error {
	data, err := redis.String(GetRedisClusterClient().GetKey(key))
	if err != nil {
		return err
	}
	if data == value {
		GetRedisClusterClient().Del(key)
		return nil
	}
	return nil
}
