package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/e"
	"github.com/spf13/viper"
)

func TestNewViperConfFail(t *testing.T) {
	if err := NewViperConf(nil); !errors.Is(err, e.ErrViperConfInvalid) {
		t.Errorf("new viper conf with nil return err not confi invalid:%v", err)
	}
	if err := NewViperConf(&ViperConf{}); !errors.Is(err, e.ErrEmptyName) {
		t.Errorf("new viper conf with nil return err not emtpy name:%v", err)
	}
}

type Zap struct {
	Level         int64  `mapstructure:"level" json:"level" yaml:"level"` //debug ,info,warn,error,panic,fatal
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"`
}

func TestViperInit(t *testing.T) {
	tc := &Zap{}
	vi := &ViperConf{
		// directory: ConfigDir,
		name: "config",
		t:    ConfigType,
		watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(tc); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			vi.SetConfigName("config")
			return nil
		},
		//
		Default: []byte(`{
"level": 0,
"stacktrace-key": "stacktrace",
"log-in-console": true}`),
	}
	defer func() {
		if err := vi.RemoveDir(); err != nil {
			t.Error(err.Error())
		}
	}()

	if vi.IsExist() {
		t.Error("config exist")
	}

	vi.Dir()
	if vi.dir != "config" {
		t.Errorf("directory want '%s' but get '%s'", "config", vi.dir)
	}

	want := Zap{
		Level:         0,
		StacktraceKey: "stacktrace",
		LogInConsole:  true,
	}

	if err := NewViperConf(vi); err != nil {
		t.Errorf("init %s's config get error: %v", str.Join(vi.name, ".", vi.t), err)
	}

	if !vi.IsExist() {
		t.Error("config not exist")
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

	dir.WriteBytes(filepath.Join(vi.getConfPath()), []byte(`{
"level": 2,
"stacktrace-key": "stacktrace1",
"log-in-console": false}`))

	want1 := Zap{
		Level:         2,
		StacktraceKey: "stacktrace1",
		LogInConsole:  false,
	}

	if err := NewViperConf(vi); err != nil {
		t.Errorf("init %s's config get error: %v", str.Join(vi.name, ".", vi.t), err)
	}

	if want1.Level != tc.Level {
		t.Errorf("want1 %+v but get %+v", want1.Level, tc.Level)
	}
	if want1.StacktraceKey != tc.StacktraceKey {
		t.Errorf("want1 %+v but get %+v", want1.StacktraceKey, tc.StacktraceKey)
	}
	if want1.LogInConsole != tc.LogInConsole {
		t.Errorf("want1 %+v but get %+v", want1.LogInConsole, tc.LogInConsole)
	}

	tc.Level = 3
	tc.StacktraceKey = "stacktrace3"
	tc.LogInConsole = true

	b, err := json.Marshal(&tc)
	if err != nil {
		t.Error(err.Error())
	}

	if err := vi.Recover(b); err != nil {
		t.Error(err.Error())
	}

	want2 := &Zap{}
	if b, err := dir.ReadBytes(vi.getConfPath()); err != nil {
		t.Error(err.Error())
	} else {
		if err := json.Unmarshal(b, want2); err != nil {
			t.Error(err.Error())
		}
	}

	if want2.Level != tc.Level {
		t.Errorf("want2 %+v but get %+v", tc.Level, want2.Level)
	}
	if want2.StacktraceKey != tc.StacktraceKey {
		t.Errorf("want2 %+v but get %+v", tc.StacktraceKey, want2.StacktraceKey)
	}
	if want2.LogInConsole != tc.LogInConsole {
		t.Errorf("want2 %+v but get %+v", tc.LogInConsole, want2.LogInConsole)
	}

	if err := vi.RemoveFile(); err != nil {
		t.Error(err.Error())
	}

	if vi.IsExist() {
		t.Error("config file exist after remove")
	}
}
