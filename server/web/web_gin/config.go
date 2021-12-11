package web_gin

import (
	"fmt"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG = Web{
	MaxSize: 1024,
	Except: Route{
		Uri:    "",
		Method: "",
	},
	System: System{
		Level:        "debug",
		Addr:         "127.0.0.1:8085",
		StaticPrefix: "/upload",
		StaticPath:   "/static/upload",
		WebPrefix:    "/admin",
		WebPath:      "./dist",
		DbType:       "mysql",
		CacheType:    "redis",
		TimeFormat:   "2006-01-02 15:04:05",
	},
	Limit: Limit{
		Disable: true,
		Limit:   0,
		Burst:   5,
	},
	Captcha: Captcha{
		KeyLong:   4,
		ImgWidth:  240,
		ImgHeight: 80,
	},
}

type Web struct {
	MaxSize int64   `mapstructure:"max-size" json:"burst" yaml:"max-size"`
	Except  Route   `mapstructure:"except" json:"except" yaml:"except"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Limit   Limit   `mapstructure:"limit" json:"limit" yaml:"limit"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
}
type Route struct {
	Uri    string `mapstructure:"uri" json:"uri" yaml:"uri"`
	Method string `mapstructure:"method" json:"method" yaml:"method"`
}

type Captcha struct {
	KeyLong   int64 `mapstructure:"key-long" json:"keyLong" yaml:"key-long"`
	ImgWidth  int64 `mapstructure:"img-width" json:"imgWidth" yaml:"img-width"`
	ImgHeight int64 `mapstructure:"img-height" json:"imgHeight" yaml:"img-height"`
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

// IsExist 配置文件是否存在
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove 删除配置文件
func Remove() error {
	err := getViperConfig().Remove()
	if err != nil {
		return fmt.Errorf("remove file %s failed %w", getViperConfig().GetConfigFileDir(), err)
	}
	return nil
}

// getViperConfig 获取初始化配置
func getViperConfig() viper_server.ViperConfig {
	maxSize := strconv.FormatInt(CONFIG.MaxSize, 10)
	keyLong := strconv.FormatInt(CONFIG.Captcha.KeyLong, 10)
	imgWidth := strconv.FormatInt(CONFIG.Captcha.ImgWidth, 10)
	imgHeight := strconv.FormatInt(CONFIG.Captcha.ImgHeight, 10)
	limit := strconv.FormatInt(int64(CONFIG.Limit.Limit), 10)
	burst := strconv.FormatInt(int64(CONFIG.Limit.Burst), 10)
	disable := strconv.FormatBool(CONFIG.Limit.Disable)
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
		max-size: ` + maxSize + `
		captcha:
		 key-long: ` + keyLong + `
		 img-width: ` + imgWidth + `
		 img-height: ` + imgHeight + `
		limit:
		 limit: ` + limit + `
		 disable: ` + disable + `
		 burst: ` + burst + `
		system:
		 level: ` + CONFIG.System.Level + `
		 addr: ` + CONFIG.System.Addr + `
		 db-type: ` + CONFIG.System.DbType + `
		 cache-type: ` + CONFIG.System.CacheType + `
		 static-path: ` + CONFIG.System.StaticPath + `
		 static-prefix: ` + CONFIG.System.StaticPrefix + `
		 time-format: ` + CONFIG.System.TimeFormat + `
		 web-prefix: ` + CONFIG.System.WebPrefix + `
		 web-path: ` + CONFIG.System.WebPath),
	}
}
