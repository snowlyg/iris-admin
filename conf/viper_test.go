package conf

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/spf13/viper"
)

type Zap struct {
	Level         int64  `mapstructure:"level" json:"level" yaml:"level"` //debug ,info,warn,error,panic,fatal
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"`
}

func TestViperInit(t *testing.T) {
	tc := &Zap{}
	config := ViperConf{
		directory: ConfigDir,
		name:      "zap", // zap => type Zap struct
		t:         ConfigType,
		watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(tc); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			vi.SetConfigName("zap")
			return nil
		},
		//
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

	defer config.RemoveDir()

	want := Zap{
		Level:         1,
		StacktraceKey: "stacktrace",
		LogInConsole:  true,
	}

	err := NewViperConf(config)
	if err != nil {
		t.Errorf("init %s's config get error: %v", str.Join(config.name, ".", config.t), err)
	}
	if want.Level != tc.Level {
		t.Errorf("want %+v but get %+v", want.Level, tc.Level)
	}
	if want.StacktraceKey != tc.StacktraceKey {
		t.Errorf("want %+v but get %+v", want.StacktraceKey, tc.StacktraceKey)
	}
	if want.LogInConsole != tc.LogInConsole {
		t.Errorf("want %+v but get %+v", want.LogInConsole, tc.LogInConsole)
	}

	dir.WriteBytes(filepath.Join(config.getConfigFilePath()), []byte(`{
"level": 2,
"stacktrace-key": "stacktrace1",
"log-in-console": false}`))

	want1 := Zap{
		Level:         2,
		StacktraceKey: "stacktrace1",
		LogInConsole:  false,
	}

	err = NewViperConf(config)
	if err != nil {
		t.Errorf("init %s's config get error: %v", str.Join(config.name, ".", config.t), err)
	}

	time.Sleep(5 * time.Second)

	if want1.Level != tc.Level {
		t.Errorf("want1 %+v but get %+v", want1.Level, tc.Level)
	}
	if want1.StacktraceKey != tc.StacktraceKey {
		t.Errorf("want1 %+v but get %+v", want1.StacktraceKey, tc.StacktraceKey)
	}
	if want1.LogInConsole != tc.LogInConsole {
		t.Errorf("want1 %+v but get %+v", want1.LogInConsole, tc.LogInConsole)
	}
}
