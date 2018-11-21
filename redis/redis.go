package redis

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"sync"
	"time"
)

type Redis struct {
	Connect   string //连接字符串
	Db        int    //数据库
	Maxidle   int    //最大空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
	Maxactive int    //最大的激活连接数，表示同时最多有N个连接
}

var (
	r           *Redis
	once        sync.Once
	redisClient *redis.Pool
)

/**
 * 返回单例实例
 * @method New
 */
func New(connect string, db int, maxidle int, maxactive int) *Redis {
	once.Do(func() { //只执行一次
		r = &Redis{Connect: connect, Db: db, Maxidle: maxidle, Maxactive: maxactive}
		setPoll()
	})
	return r
}

/**
 * 公共方法
 */
/**
 * 设置连接池
 * @method setPoll
 */
func setPoll() {
	redisClient = &redis.Pool{
		MaxIdle:     r.Maxidle,
		MaxActive:   r.Maxactive,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) { //建立连接
			log.Printf(r.Connect)
			c, err := redis.Dial("tcp", r.Connect)
			if err != nil {
				panic(err)
			}
			c.Do("SELECT", r.Db)
			return c, nil
		},
	}
}

/**
 * 执行基本命令
 * @method func
 * @param  {[type]} n *Neo4j        [description]
 * @return {[type]}   [description]
 */
func (n *Redis) Do(cmd string, args ...interface{}) (interface{}, error) {
	conn := redisClient.Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

/**
 * 设置键值对, ex单位是秒
 * @method func
 * @param  {[type]} n *Redis        [description]
 * @return {[type]}   [description]
 */
func (n *Redis) SetString(key string, value string, ex string) (interface{}, error) {
	conn := redisClient.Get()
	defer conn.Close()
	return conn.Do("SET", key, value, "EX", ex)
}

/**
 * 获取键的值
 * @method func
 * @param  {[type]} n *Redis        [description]
 * @return {[type]}   [description]
 */
func (n *Redis) GetString(key string) (string, error) {
	conn := redisClient.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do("GET", key))
	return value, err
}
