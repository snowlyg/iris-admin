package conf

import (
	"bytes"
	"errors"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/snowlyg/iris-admin/e"
	"github.com/spf13/viper"
)

func TestSetDefaultAddrAndTimeFormat(t *testing.T) {
	clearConfigEnv(t)
	dc := &Conf{}
	addr := ""
	if dc.System.Addr != addr {
		t.Errorf("config system addr want '%s' but get '%s'", addr, dc.System.Addr)
	}
	timeFormat := ""
	if dc.System.TimeFormat != timeFormat {
		t.Errorf("config system time format want '%s' but get '%s'", timeFormat, dc.System.TimeFormat)
	}
	dc.SetDefaultAddrAndTimeFormat()
	addr = "127.0.0.1:8080"
	if dc.System.Addr != addr {
		t.Errorf("config system addr want '%s' but get '%s'", addr, dc.System.Addr)
	}
	timeFormat = "2006-01-02 15:04:05"
	if dc.System.TimeFormat != timeFormat {
		t.Errorf("config system time format want '%s' but get '%s'", timeFormat, dc.System.TimeFormat)
	}

	c := NewConf()
	c.configDir = t.TempDir()
	if c.IsExist() {
		t.Error("config exist before init")
	}
	addr = "127.0.0.1:8080"
	if c.System.Addr != addr {
		t.Errorf("config system addr want '%s' but get '%s'", addr, c.System.Addr)
	}
	timeFormat = "2006-01-02 15:04:05"
	if c.System.TimeFormat != timeFormat {
		t.Errorf("config system time format want '%s' but get '%s'", timeFormat, c.System.TimeFormat)
	}

	if err := c.Recover(); err != nil {
		t.Error(err.Error())
	}
	if !c.IsExist() {
		t.Error("config not exist after recover")
	}
	c.RemoveFile()
	if c.IsExist() {
		t.Error("config exist after remove")
	}
}

func clearConfigEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		mysqlAddrKey,
		mysqlPwdKey,
		mysqlNameKey,
		mysqlUserKey,
		mysqlDbNameKey,
		webAddrKey,
	} {
		t.Setenv(key, "")
	}
}

func TestApplyEnvAndWarningsHandleNilConfiguration(t *testing.T) {
	var conf *Conf
	conf.applyEnv()
	conf.logWarnings()

	conf = &Conf{}
	conf.applyEnv()
	conf.logWarnings()
}

func TestLogWarningsReportsIncompleteConfiguration(t *testing.T) {
	var output bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&output)
	defer log.SetOutput(previous)

	conf := newDefaultConf()
	conf.Mysql.Path = ""
	conf.logWarnings()
	if !strings.Contains(output.String(), mysqlAddrKey) {
		t.Fatalf("warning = %q, want environment key", output.String())
	}
}

func TestRecoverPropagatesRbacFailure(t *testing.T) {
	var conf *Conf
	if err := conf.Recover(); err == nil {
		t.Fatal("Recover ignored invalid configuration")
	}
}

func TestFileMaxSizeYAMLTag(t *testing.T) {
	field, ok := reflect.TypeOf(Conf{}).FieldByName("FileMaxSize")
	if !ok {
		t.Fatal("FileMaxSize field not found")
	}
	if got := field.Tag.Get("yaml"); got != "file-max-size" {
		t.Fatalf("FileMaxSize yaml tag = %q, want file-max-size", got)
	}
}

func TestConfigWatchReportsUnmarshalFailure(t *testing.T) {
	conf := newDefaultConf()
	watch := conf.getViperConfig().watch
	vi := viper.New()
	vi.Set("mysql", "invalid")
	if err := watch(vi); err == nil {
		t.Fatal("watch accepted invalid mysql configuration")
	}
}

func TestRecoverReturnsEnsureRbacError(t *testing.T) {
	conf := newDefaultConf()
	conf.configDir = "\x00"
	if err := conf.Recover(); err == nil {
		t.Fatal("Recover ignored RBAC path error")
	}
}

func TestNewConfValidation(t *testing.T) {
	clearConfigEnv(t)
	conf := NewConf()
	if err := conf.Validate(); err != nil {
		t.Fatalf("NewConf validation: %v", err)
	}
	if err := (*Conf)(nil).Validate(); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("nil configuration error = %v, want %v", err, e.ErrConfigInvalid)
	}
}
