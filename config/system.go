package config

type System struct {
	Level     string `mapstructure:"level" json:"level" yaml:"level"` // debug,release,test
	Addr      string `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType    string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	CacheType string `mapstructure:"cache-type" json:"cacheType" yaml:"cache-type"`
}
