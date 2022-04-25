package zap_server

import (
	"fmt"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG = Zap{
	Level:         "debug",
	Format:        "console",
	Prefix:        "[IRIS-ADMIN]",
	Director:      "logs",
	LinkName:      "latest_log",
	ShowLine:      true,
	EncodeLevel:   "LowercaseColorLevelEncoder",
	StacktraceKey: "stacktrace",
	LogInConsole:  false,
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"` //debug ,info,warn,error,panic,fatal
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`
	LinkName      string `mapstructure:"link-name" json:"link-name" yaml:"link-name"`
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"`
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

// getViperConfig 获取初始化配置
func getViperConfig() viper_server.ViperConfig {
	configName := "zap"
	showLine := strconv.FormatBool(CONFIG.ShowLine)
	logInConsole := strconv.FormatBool(CONFIG.LogInConsole)
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
{
	"level": "` + CONFIG.Level + `",
	"format": "` + CONFIG.Format + `",
	"prefix": "` + CONFIG.Prefix + `",
	"director": "` + CONFIG.Director + `",
	"link-name": "` + CONFIG.LinkName + `",
	"show-line": ` + showLine + `,
	"encode-level": "` + CONFIG.EncodeLevel + `",
	"stacktrace-key": "` + CONFIG.StacktraceKey + `",
	"log-in-console": ` + logInConsole + `
}`),
	}
}
