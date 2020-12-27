package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"time"
)

type RedisCluster struct {
	*redisc.Cluster
}

var rcClient *RedisCluster

func createPool(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr, opts...)
			if err != nil {
				fmt.Println(fmt.Sprintf("dial redis error:%+v", err))
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				fmt.Println(fmt.Sprintf("从 redis 连接池取出连接无效：%+v", err))
			}
			return err
		},
		Wait: true,
	}, nil
}

func InitRedisCluster(addrs []string, password string) {
	rc := &redisc.Cluster{
		StartupNodes: addrs,
		DialOptions: []redis.DialOption{
			redis.DialConnectTimeout(3 * time.Second),
			redis.DialPassword(password),
		},
		CreatePool: createPool,
	}
	rc.Refresh()
	rcClient = &RedisCluster{
		Cluster: rc,
	}
}

func GetRedisClusterClient() *RedisCluster {
	return rcClient
}

func (rc *RedisCluster) Do(cmd string, args ...interface{}) (interface{}, error) {
	conn := rc.Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

func (rc *RedisCluster) Send(cmd string, args ...interface{}) error {
	conn := rc.Get()
	defer conn.Close()
	return conn.Send(cmd, args...)
}

func (rc *RedisCluster) Close() {
	return
}

func (rc *RedisCluster) GetKey(key string) (interface{}, error) {
	return rc.Do("GET", key)
}

func (rc *RedisCluster) Set(key string, value interface{}, ttl ...time.Duration) (interface{}, error) {
	reply, err := rc.Do("SET", key, value)
	if len(ttl) <= 0 {
		return reply, err
	}
	if _, err := rc.Expire(key, int(ttl[0].Seconds())); err != nil {
		fmt.Println(fmt.Sprintf("%s.EXPIRE.err : %s", key, err.Error()))
	}
	return reply, err
}

func (rc *RedisCluster) SetNX(key string, value interface{}, expireSeconds int) bool {
	n, err := rc.Do("SETNX", key, value)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s.SetNX.err : %s", key, err.Error()))
		return false
	}
	if n == 1 {
		rc.Expire(key, expireSeconds)
		return true
	}
	return false
}

func (rc *RedisCluster) Del(keys ...interface{}) (int, error) {
	n := 0
	t := 0
	var err error
	for _, v := range keys {
		t, err = redis.Int(rc.Do("DEL", v))
		n += t
	}
	return n, err
}

func (rc *RedisCluster) Exists(key string) bool {
	n, err := redis.Int(rc.Do("EXISTS", key))
	if err != nil {
		fmt.Println(fmt.Sprintf("%s.Exists.err : %s", key, err.Error()))
		return false
	}
	if n == 1 {
		return true
	}

	return false
}

func (rc *RedisCluster) Expire(key string, seconds int) (interface{}, error) {
	if seconds <= 0 {
		fmt.Println(fmt.Sprintf("redis expire invalid second: %d", seconds))
		return nil, nil
	}
	return rc.Do("EXPIRE", key, seconds)
}

func (rc *RedisCluster) LPush(key string, values ...interface{}) (interface{}, error) {
	data := make([]interface{}, 0, len(values)+1)
	data = append(data, key)
	data = append(data, values...)
	return rc.Do("LPUSH", data)
}

func (rc *RedisCluster) RPush(key string, values ...interface{}) (interface{}, error) {
	data := make([]interface{}, 0, len(values)+1)
	data = append(data, key)
	data = append(data, values...)
	return rc.Do("RPUSH", data)
}

func (rc *RedisCluster) LLen(key string) int {
	length, _ := redis.Int(rc.Do("LLEN", key))
	return length
}

func (rc *RedisCluster) LTrim(key string, start, end int) (interface{}, error) {
	return rc.Do("LTRIM", key, start, end)
}

func (rc *RedisCluster) LRange(key string, start, end int) (interface{}, error) {
	return rc.Do("LRANGE", key, start, end)
}

func (rc *RedisCluster) HGetAll(key string) (interface{}, error) {
	return rc.Do("HGETALL", key)
}

func (rc *RedisCluster) HMSet(key string, values ...interface{}) (interface{}, error) {
	data := make([]interface{}, 0, len(values)+1)
	data = append(data, key)
	data = append(data, values...)
	return rc.Do("HMSET", data...)
}

func (rc *RedisCluster) HIncrBy(key string, field string, incr int64) (interface{}, error) {
	return rc.Do("HINCRBY", key, field, incr)
}

func (rc *RedisCluster) Sadd(key string, members ...interface{}) (interface{}, error) {
	data := make([]interface{}, 0, len(members)+1)
	data = append(data, key)
	data = append(data, members...)
	return rc.Do("SADD", data...)
}

func (rc *RedisCluster) Scard(key string) (interface{}, error) {
	return rc.Do("SCARD", key)
}

func (rc *RedisCluster) Members(key string) (interface{}, error) {
	return rc.Do("SMEMBERS", key)
}

// LoadRedisHashToStruct 从 redis 加载数据
func (rc *RedisCluster) LoadRedisHashToStruct(sKey string, pst interface{}) error {
	vals, err := redis.Values(rc.HGetAll(sKey))
	if err != nil {
		return err
	}
	err = redis.ScanStruct(vals, pst)
	if err != nil {
		return err
	}
	return nil
}
