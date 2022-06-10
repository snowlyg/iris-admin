package viper_server

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"github.com/spf13/viper"
)

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

func TestViperInit(t *testing.T) {
	tc := &Zap{}
	config := ViperConfig{
		Directory: g.ConfigDir,
		Name:      "zap", // 名称需要和结构体名称对应 zap => type Zap struct
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(tc); err != nil {
				return fmt.Errorf("反序列化错误: %v", err)
			}
			vi.SetConfigName("zap")
			return nil
		},
		// 注意:设置默认配置值的时候,前面不能有空格等其他符号.必须紧贴左侧.
		Default: []byte(`{
"level": "info",
"format": "console",
"prefix": "[OP-ONLINE]",
"director": "log",
"link-name": "latest_log",
"show-line": true,
"encode-level": "LowercaseColorLevelEncoder",
"stacktrace-key": "stacktrace",
"log-in-console": true}`),
	}

	defer config.Remove()

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

	err := Init(config)
	if err != nil {
		t.Errorf("初始化 %s 的配置返回错误: %v", str.Join(config.Name, ".", config.Type), err)
	}
	if want.Level != tc.Level {
		t.Errorf("want %+v but get %+v", want.Level, tc.Level)
	}
	if want.Format != tc.Format {
		t.Errorf("want %+v but get %+v", want.Format, tc.Format)
	}
	if want.Prefix != tc.Prefix {
		t.Errorf("want %+v but get %+v", want.Prefix, tc.Prefix)
	}
	if want.Director != tc.Director {
		t.Errorf("want %+v but get %+v", want.Director, tc.Director)
	}
	if want.LinkName != tc.LinkName {
		t.Errorf("want %+v but get %+v", want.LinkName, tc.LinkName)
	}
	if want.ShowLine != tc.ShowLine {
		t.Errorf("want %+v but get %+v", want.ShowLine, tc.ShowLine)
	}
	if want.EncodeLevel != tc.EncodeLevel {
		t.Errorf("want %+v but get %+v", want.EncodeLevel, tc.EncodeLevel)
	}
	if want.StacktraceKey != tc.StacktraceKey {
		t.Errorf("want %+v but get %+v", want.StacktraceKey, tc.StacktraceKey)
	}
	if want.LogInConsole != tc.LogInConsole {
		t.Errorf("want %+v but get %+v", want.LogInConsole, tc.LogInConsole)
	}

	dir.WriteBytes(filepath.Join(config.getConfigFilePath()), []byte(`{
"level": "info1",
"format": "console1",
"prefix": "[OP-ONLINE]1",
"director": "log1",
"link-name": "latest_log1",
"show-line": false,
"encode-level": "LowercaseColorLevelEncoder1",
"stacktrace-key": "stacktrace1",
"log-in-console": false}`))

	want1 := Zap{
		Level:         "info1",
		Format:        "console1",
		Prefix:        "[OP-ONLINE]1",
		Director:      "log1",
		LinkName:      "latest_log1",
		ShowLine:      false,
		EncodeLevel:   "LowercaseColorLevelEncoder1",
		StacktraceKey: "stacktrace1",
		LogInConsole:  false,
	}

	err = Init(config)
	if err != nil {
		t.Errorf("初始化 %s 的配置返回错误: %v", str.Join(config.Name, ".", config.Type), err)
	}

	time.Sleep(5 * time.Second)

	if want1.Level != tc.Level {
		t.Errorf("want1 %+v but get %+v", want1.Level, tc.Level)
	}
	if want1.Format != tc.Format {
		t.Errorf("want1 %+v but get %+v", want1.Format, tc.Format)
	}
	if want1.Prefix != tc.Prefix {
		t.Errorf("want1 %+v but get %+v", want1.Prefix, tc.Prefix)
	}
	if want1.Director != tc.Director {
		t.Errorf("want1 %+v but get %+v", want1.Director, tc.Director)
	}
	if want1.LinkName != tc.LinkName {
		t.Errorf("want1 %+v but get %+v", want1.LinkName, tc.LinkName)
	}
	if want1.ShowLine != tc.ShowLine {
		t.Errorf("want1 %+v but get %+v", want1.ShowLine, tc.ShowLine)
	}
	if want1.EncodeLevel != tc.EncodeLevel {
		t.Errorf("want1 %+v but get %+v", want1.EncodeLevel, tc.EncodeLevel)
	}
	if want1.StacktraceKey != tc.StacktraceKey {
		t.Errorf("want1 %+v but get %+v", want1.StacktraceKey, tc.StacktraceKey)
	}
	if want1.LogInConsole != tc.LogInConsole {
		t.Errorf("want1 %+v but get %+v", want1.LogInConsole, tc.LogInConsole)
	}
}
