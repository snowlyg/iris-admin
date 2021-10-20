package viper_server

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"github.com/spf13/viper"
)

var tc Zap

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"` //debug ,info,warn,error,panic,fatal
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`
	LinkName      string `mapstructure:"link-name" json:"linkName" yaml:"link-name"`
	ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"show-line"`
	EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`
}

func TestViperInit(t *testing.T) {
	config := ViperConfig{
		Directory: g.ConfigDir,
		Name:      "zap", // 名称需要和结构体名称对应 zap => type Zap struct
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&tc); err != nil {
				return fmt.Errorf("反序列化错误: %v", err)
			}
			vi.SetConfigName("zap")
			// 监控配置文件变化
			vi.WatchConfig()
			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("配置发生变化:", e.Name)
				if err := vi.Unmarshal(&tc); err != nil {
					fmt.Printf("反序列化错误: %v \n", err)
				}
			})
			return nil
		},
		// 注意:设置默认配置值的时候,前面不能有空格等其他符号.必须紧贴左侧.
		Default: []byte(`
level: info
format: console
prefix: '[OP-ONLINE]'
director: log
link-name: latest_log
show-line: true
encode-level: LowercaseColorLevelEncoder
stacktrace-key: stacktrace
log-in-console: true`),
	}
	want := Zap{
		Level:         "info",
		Format:        "console",
		Prefix:        "[OP-ONLINE]",
		Director:      "log",
		LinkName:      "latest_log",
		ShowLine:      true,
		EncodeLevel:   "LowercaseColorLevelEncoder",
		StacktraceKey: "stacktrace",
		LogInConsole:  true,
	}
	t.Run("test viper init func", func(t *testing.T) {
		err := Init(config)
		if err != nil {
			t.Errorf("初始化 %s 的配置返回错误: %v", str.Join(config.Name, ".", config.Type), err)
		}
		if !reflect.DeepEqual(want, tc) {
			t.Errorf("test viper init want %+v but get %+v", want, tc)
		}
	})
}
