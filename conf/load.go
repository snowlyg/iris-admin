package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/e"
)

const configName = "iris_admin.json"

var absConfigPath = filepath.Abs

// LoadConf loads defaults, overlays the default JSON file when present, and
// applies environment variables last.
func LoadConf() (*Conf, error) {
	return LoadConfFile(filepath.Join(defaultConfigDirectory(), configName))
}

// LoadConfFile loads defaults, overlays filePath when present, and applies
// environment variables last. A missing file is not an error.
func LoadConfFile(filePath string) (*Conf, error) {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return nil, fmt.Errorf("config path is empty: %w", e.ErrConfigInvalid)
	}
	absolutePath, err := absConfigPath(filePath)
	if err != nil {
		return nil, fmt.Errorf("resolve config path %q: %w", filePath, err)
	}
	filePath = absolutePath

	conf := newDefaultConf()
	conf.configDir = filepath.Dir(filePath)
	data, err := os.ReadFile(filePath)
	switch {
	case err == nil:
		if err := json.Unmarshal(data, conf); err != nil {
			return nil, fmt.Errorf("decode config %q: %w", filePath, err)
		}
	case !os.IsNotExist(err):
		return nil, fmt.Errorf("read config %q: %w", filePath, err)
	}

	conf.applyEnv()
	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf("validate config %q: %w", filePath, err)
	}
	conf.logWarnings()
	return conf, nil
}

// Validate checks configuration required during service startup.
func (conf *Conf) Validate() error {
	if conf == nil {
		return fmt.Errorf("config is nil: %w", e.ErrConfigInvalid)
	}
	if strings.TrimSpace(conf.System.Addr) == "" {
		return fmt.Errorf("system address is empty: %w", e.ErrConfigInvalid)
	}
	if conf.Mysql == nil {
		return fmt.Errorf("mysql config is nil: %w", e.ErrConfigInvalid)
	}
	if strings.TrimSpace(conf.Mysql.Path) == "" {
		return fmt.Errorf("mysql address is empty: %w", e.ErrConfigInvalid)
	}
	if strings.TrimSpace(conf.Mysql.DbName) == "" {
		return fmt.Errorf("mysql database name is empty: %w", e.ErrConfigInvalid)
	}
	if strings.TrimSpace(conf.Mysql.Username) == "" {
		return fmt.Errorf("mysql username is empty: %w", e.ErrConfigInvalid)
	}
	return nil
}

func defaultConfigDirectory() string {
	return filepath.Join(dir.GetCurrentAbPath(), ConfigDir)
}

func (conf *Conf) configDirectory() string {
	if conf != nil && conf.configDir != "" {
		return filepath.Clean(conf.configDir)
	}
	return defaultConfigDirectory()
}
