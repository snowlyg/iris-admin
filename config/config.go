package config

import (
	"github.com/BurntSushi/toml"
	"sync"
)

type Config struct {
	App     app
	Mysql   mysql
	Mongodb mongodb
	Redis   redis
}

type app struct {
	Addr string
}

type mysql struct {
	Connect string
	MaxIdle int
	MaxOpen int
}

type mongodb struct {
	Connect string
}

type redis struct {
	Connect   string
	DB        int
	MaxIdle   int
	MaxActive int
}

var (
	c    *Config
	once sync.Once
)

/**
 * 返回单例实例
 * @method New
 */
func New() *Config {
	once.Do(func() { //只执行一次
		if _, err := toml.DecodeFile("config.toml", &c); err != nil {
			panic(err.Error())
		}
	})
	return c
}
