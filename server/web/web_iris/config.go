package web_iris

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG Web

type Web struct {
	MaxSize int64   `mapstructure:"max-size" json:"burst" yaml:"max-size"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Limit   Limit   `mapstructure:"limit" json:"limit" yaml:"limit"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
}

type Captcha struct {
	KeyLong   int `mapstructure:"key-long" json:"keyLong" yaml:"key-long"`
	ImgWidth  int `mapstructure:"img-width" json:"imgWidth" yaml:"img-width"`
	ImgHeight int `mapstructure:"img-height" json:"imgHeight" yaml:"img-height"`
}

type Limit struct {
	Disable bool    `mapstructure:"disable" json:"disable" yaml:"disable"`
	Limit   float64 `mapstructure:"limit" json:"limit" yaml:"limit"`
	Burst   int     `mapstructure:"burst" json:"burst" yaml:"burst"`
}

type System struct {
	Level        string `mapstructure:"level" json:"level" yaml:"level"` // debug,release,test
	Addr         string `mapstructure:"addr" json:"addr" yaml:"addr"`
	StaticPrefix string `mapstructure:"static-prefix" json:"staticPrefix" yaml:"static-prefix"`
	StaticPath   string `mapstructure:"static-path" json:"staticPath" yaml:"static-path"`
	WebPrefix    string `mapstructure:"web-prefix" json:"webPPrefix" yaml:"web-prefix"`
	WebPath      string `mapstructure:"web-path" json:"webPath" yaml:"web-path"`
	DbType       string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	CacheType    string `mapstructure:"cache-type" json:"cacheType" yaml:"cache-type"`
	TimeFormat   string `mapstructure:"time-format" json:"timeFormat" yaml:"time-format"`
}

// getViperConfig 获取初始化配置
func getViperConfig() viper_server.ViperConfig {
	configName := "web"
	return viper_server.ViperConfig{
		Directory: g.ConfigDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("反序列化错误: %v", err)
			}
			// 监控配置文件变化
			vi.SetConfigName(configName)
			vi.WatchConfig()
			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("配置发生变化:", e.Name)
				if err := vi.Unmarshal(&CONFIG); err != nil {
					fmt.Printf("反序列化错误: %v \n", err)
				}
			})
			return nil
		},
		// 注意:设置默认配置值的时候,前面不能有空格等其他符号.必须紧贴左侧.
		Default: []byte(`
max-size: 1024
captcha:
 key-long: 6
 img-width: 240
 img-height: 80
limit:
 limit: false
 limit: 0
 burst: 5
system:
 level: debug
 addr: 127.0.0.1:8085
 db-type: mysql
 cache-type: local
 static-path: /static/upload
 static-prefix: /upload
 time-format: "2006-01-02 15:04:05"
 web-prefix: /
 web-path: ./dist`),
	}
}
