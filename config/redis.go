package config

type Redis struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	PoolSize int    `mapstructure:"pool-size" json:"poolSize" yaml:"pool-size"`
}
