package conf

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/snowlyg/iris-admin/e"
)

func TestLoadConfFilePrecedence(t *testing.T) {
	clearConfigEnv(t)
	path := filepath.Join(t.TempDir(), configName)
	data := []byte(`{
		"locale": "en",
		"system": {"addr": "file-host:8080"},
		"mysql": {
			"path": "file-db:3306",
			"db-name": "file-name",
			"username": "file-user",
			"max-open-conns": 12
		}
	}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv(webAddrKey, "env-host:9090")
	t.Setenv(mysqlAddrKey, "env-db:3307")
	t.Setenv(mysqlUserKey, "env-user")
	t.Setenv(mysqlDbNameKey, "env-name")
	t.Setenv(mysqlPwdKey, "env-password")

	conf, err := LoadConfFile(path)
	if err != nil {
		t.Fatalf("LoadConfFile: %v", err)
	}
	if conf.Locale != "en" {
		t.Fatalf("locale = %q, want file value", conf.Locale)
	}
	if conf.System.Addr != "env-host:9090" {
		t.Fatalf("system address = %q, want environment value", conf.System.Addr)
	}
	if conf.Mysql.Path != "env-db:3307" ||
		conf.Mysql.Username != "env-user" ||
		conf.Mysql.DbName != "env-name" ||
		conf.Mysql.Password != "env-password" {
		t.Fatalf("mysql environment overrides not applied: %+v", conf.Mysql)
	}
	if conf.Mysql.MaxOpenConns != 12 {
		t.Fatalf("max open connections = %d, want file value 12", conf.Mysql.MaxOpenConns)
	}
	if conf.Mysql.Config == "" {
		t.Fatal("partial file overlay removed default mysql config")
	}
	if conf.configDirectory() != filepath.Dir(path) {
		t.Fatalf("config directory = %q, want %q", conf.configDirectory(), filepath.Dir(path))
	}
}

func TestLoadConfFileMissingUsesDefaultsAndLegacyEnv(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv(mysqlNameKey, "legacy-user")
	t.Setenv(mysqlPwdKey, "password")

	conf, err := LoadConfFile(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("LoadConfFile missing file: %v", err)
	}
	if conf.Locale != "zh" || conf.System.Addr != "127.0.0.1:8080" {
		t.Fatalf("defaults not preserved: %+v", conf)
	}
	if conf.Mysql.Username != "legacy-user" {
		t.Fatalf("legacy username override = %q, want legacy-user", conf.Mysql.Username)
	}
}

func TestLoadConfFileRejectsInvalidInput(t *testing.T) {
	clearConfigEnv(t)
	if _, err := LoadConfFile(""); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("empty path error = %v, want %v", err, e.ErrConfigInvalid)
	}

	dir := t.TempDir()
	invalidJSON := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(invalidJSON, []byte(`{"mysql":`), 0o600); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}
	if _, err := LoadConfFile(invalidJSON); err == nil {
		t.Fatal("invalid JSON was accepted")
	}

	nilMysql := filepath.Join(dir, "nil-mysql.json")
	if err := os.WriteFile(nilMysql, []byte(`{"mysql": null}`), 0o600); err != nil {
		t.Fatalf("write nil mysql config: %v", err)
	}
	if _, err := LoadConfFile(nilMysql); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("nil mysql error = %v, want %v", err, e.ErrConfigInvalid)
	}

	if _, err := LoadConfFile(dir); err == nil {
		t.Fatal("directory config path was accepted")
	}
}

func TestLoadConfUsesDefaultPath(t *testing.T) {
	clearConfigEnv(t)
	conf, err := LoadConf()
	if err != nil {
		t.Fatalf("LoadConf: %v", err)
	}
	if conf == nil || conf.configDirectory() != defaultConfigDirectory() {
		t.Fatalf("default config directory = %q, want %q", conf.configDirectory(), defaultConfigDirectory())
	}
}

func TestLoadConfFileReportsPathResolutionFailure(t *testing.T) {
	previous := absConfigPath
	absConfigPath = func(string) (string, error) {
		return "", errors.New("resolve failed")
	}
	defer func() {
		absConfigPath = previous
	}()

	if _, err := LoadConfFile("config.json"); err == nil {
		t.Fatal("path resolution failure was ignored")
	}
}

func TestValidateConfiguration(t *testing.T) {
	valid := newDefaultConf()
	tests := []struct {
		name string
		conf *Conf
	}{
		{name: "nil", conf: nil},
		{name: "empty address", conf: mutateConf(valid, func(c *Conf) { c.System.Addr = " " })},
		{name: "nil mysql", conf: mutateConf(valid, func(c *Conf) { c.Mysql = nil })},
		{name: "empty mysql address", conf: mutateConf(valid, func(c *Conf) { c.Mysql.Path = " " })},
		{name: "empty database", conf: mutateConf(valid, func(c *Conf) { c.Mysql.DbName = " " })},
		{name: "empty username", conf: mutateConf(valid, func(c *Conf) { c.Mysql.Username = " " })},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.conf.Validate(); !errors.Is(err, e.ErrConfigInvalid) {
				t.Fatalf("Validate error = %v, want %v", err, e.ErrConfigInvalid)
			}
		})
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid configuration error: %v", err)
	}
}

func TestConfigDirectoryFallback(t *testing.T) {
	var nilConf *Conf
	if got := nilConf.configDirectory(); got != defaultConfigDirectory() {
		t.Fatalf("nil config directory = %q, want %q", got, defaultConfigDirectory())
	}
}

func mutateConf(base *Conf, mutate func(*Conf)) *Conf {
	copy := *base
	if base.Mysql != nil {
		mysqlCopy := *base.Mysql
		copy.Mysql = &mysqlCopy
	}
	mutate(&copy)
	return &copy
}

func TestLoadConfFileWrapsValidationContext(t *testing.T) {
	clearConfigEnv(t)
	path := filepath.Join(t.TempDir(), "invalid-config.json")
	if err := os.WriteFile(path, []byte(`{"system":{"addr":" "}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfFile(path); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("validation error = %v, want wrapped %v", err, e.ErrConfigInvalid)
	}
}
