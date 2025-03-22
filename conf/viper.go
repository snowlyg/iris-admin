package conf

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/e"
	"github.com/spf13/viper"
)

type ViperConf struct {
	directory string
	name      string
	t         string
	Default   []byte
	watch     func(*viper.Viper) error
}

// getConfigFilePath
func (vc ViperConf) getConfigFilePath() string {
	return filepath.Join(dir.GetCurrentAbPath(), vc.directory, str.Join(vc.name, ".", vc.t))
}

// getConfigFileDir
func (vc ViperConf) getConfigFileDir() string {
	if vc.directory == "" {
		return "config"
	}
	return vc.directory
}

// IsFileExist
func (vc ViperConf) IsFileExist() bool {
	return dir.IsExist(vc.getConfigFilePath())
}

// RemoveFile remove config file
func (vc ViperConf) RemoveFile() error {
	return dir.Remove(vc.getConfigFilePath())
}

// RemoveDir remove config dir
func (vc ViperConf) RemoveDir() error {
	return os.RemoveAll(filepath.Dir(vc.getConfigFilePath()))
}

// Recover recover config file content
func (vc ViperConf) Recover(b []byte) error {
	_, err := dir.WriteBytes(vc.getConfigFilePath(), b)
	return err
}

// NewViperConf
func NewViperConf(vc ViperConf) error {
	if vc.name == "" {
		return e.ErrEmptyName
	}
	if vc.t == "" {
		vc.t = "yaml"
	}

	vc.directory = vc.getConfigFileDir()
	filePath := vc.getConfigFilePath()

	vi := viper.New()
	vi.SetConfigName(vc.name)
	vi.SetConfigType(vc.t)
	vi.AddConfigPath(vc.directory)
	isExist := dir.IsExist(filePath)
	if !isExist {
		if vc.directory != "./" {
			err := dir.InsureDir(filepath.Dir(filePath))
			if err != nil {
				return fmt.Errorf("create dir %s fail : %v", filePath, err)
			}
		}
		// ReadConfig
		if err := vi.ReadConfig(bytes.NewBuffer(vc.Default)); err != nil {
			return fmt.Errorf("read default config fail : %w ", err)
		}
		// WriteConfigAs
		if err := vi.WriteConfigAs(filePath); err != nil {
			return fmt.Errorf("write config to path fail: %w ", err)
		}
	} else {
		vi.SetConfigFile(filePath)
		if err := vi.ReadInConfig(); err != nil {
			return fmt.Errorf("read config fail: %w ", err)
		}
	}
	if err := vc.watch(vi); err != nil {
		return fmt.Errorf("watch config fail: %w ", err)
	}
	return nil
}
