package mongodb

import (
	"fmt"
	"time"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG = MongoDB{
	DB:      "mongo_test",
	Timeout: 10,
	Addr:    "localhost:27017",
}

type MongoDB struct {
	Timeout time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	DB      string        `mapstructure:"db" json:"db" yaml:"db"`
	Addr    string        `mapstructure:"addr" json:"addr" yaml:"addr"`
}

func (md *MongoDB) GetApplyURI() string {
	return str.Join("mongodb://", md.Addr)
}

// IsExist 配置文件是否存在
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove 删除配置文件
func Remove() error {
	err := getViperConfig().Remove()
	if err != nil {
		return err
	}
	return nil
}

// init 初始化配置
func init() {
	viper_server.Init(getViperConfig())
}

// getViperConfig 获取初始化配置
func getViperConfig() viper_server.ViperConfig {
	configName := "mongo"
	return viper_server.ViperConfig{
		Debug:     true,
		Directory: g.ConfigDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("反序列化错误: %v", err)
			}
			// 监控配置文件变化
			vi.SetConfigName(configName)
			return nil
		},
		// 注意:设置默认配置值的时候,前面不能有空格等其他符号.必须紧贴左侧.
		Default: []byte(`
{
	"timeout": "` + CONFIG.Timeout.String() + `",
	"db": "` + CONFIG.DB + `",
	"addr": "` + CONFIG.Addr + `"
}`),
	}
}
