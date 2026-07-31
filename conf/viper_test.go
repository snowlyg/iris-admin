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
	if err := NewViperConf(&ViperConf{}); !errors.Is(err, e.ErrConfigNameEmpty) {
		t.Errorf("new viper conf with nil return err not emtpy name:%v", err)
	}
	var conf *ViperConf
	if conf.getConfPath() != "" {
		t.Fatal("nil ViperConf returned a config path")
	}
	if conf.IsExist() {
		t.Fatal("nil ViperConf exists")
	}
	if err := conf.RemoveFile(); !errors.Is(err, e.ErrViperConfInvalid) {
		t.Fatalf("nil RemoveFile error = %v", err)
	}
	if err := conf.RemoveDir(); !errors.Is(err, e.ErrViperConfInvalid) {
		t.Fatalf("nil RemoveDir error = %v", err)
	}
	if err := conf.Recover(nil); !errors.Is(err, e.ErrViperConfInvalid) {
		t.Fatalf("nil Recover error = %v", err)
	}
}

type Zap struct {
	Level         int64  `mapstructure:"level" json:"level" yaml:"level"` //debug ,info,warn,error,panic,fatal
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"`
}

func TestViperInit(t *testing.T) {
	tc := &Zap{}
	configDir := t.TempDir()
	vi := &ViperConf{
		dir:  configDir,
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
	if vi.IsExist() {
		t.Error("config exist")
	}

	vi.Dir()
	if vi.dir != configDir {
		t.Errorf("directory want '%s' but get '%s'", configDir, vi.dir)
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

func TestViperConfResolvesRelativeAndAbsoluteDirectories(t *testing.T) {
	relative := &ViperConf{dir: "relative", name: "config", t: "json"}
	if !filepath.IsAbs(relative.resolvedDir()) {
		t.Fatalf("relative resolved directory = %q, want absolute", relative.resolvedDir())
	}

	absoluteDir := t.TempDir()
	absolute := &ViperConf{dir: absoluteDir, name: "config", t: "json"}
	if got := absolute.resolvedDir(); got != filepath.Clean(absoluteDir) {
		t.Fatalf("absolute resolved directory = %q, want %q", got, absoluteDir)
	}
}

func TestViperConfRefusesPathTraversalRemoval(t *testing.T) {
	root := t.TempDir()
	conf := &ViperConf{dir: root, name: "../outside", t: "json"}
	if err := conf.RemoveFile(); err == nil {
		t.Fatal("RemoveFile accepted path traversal")
	}
	if err := conf.RemoveDir(); err == nil {
		t.Fatal("RemoveDir accepted path traversal")
	}
}

func TestViperConfRemovesOwnedDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "owned")
	conf := &ViperConf{dir: root, name: "config", t: "json"}
	if err := conf.Recover([]byte(`{}`)); err != nil {
		t.Fatalf("Recover: %v", err)
	}
	if err := conf.RemoveDir(); err != nil {
		t.Fatalf("RemoveDir: %v", err)
	}
	if dir.IsExist(root) {
		t.Fatalf("owned directory %q still exists", root)
	}
}
