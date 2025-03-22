package viper_server

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/spf13/viper"
)

var (
	ErrEmptyName = errors.New("config'name can't be empty value")
)

type ViperConfig struct {
	Directory string
	Name      string
	Type      string
	Default   []byte
	Watch     func(*viper.Viper) error
}

// getConfigFilePath
func (vc ViperConfig) getConfigFilePath() string {
	return filepath.Join(dir.GetCurrentAbPath(), vc.Directory, str.Join(vc.Name, ".", vc.Type))
}

// getConfigFileDir
func (vc ViperConfig) getConfigFileDir() string {
	if vc.Directory == "" {
		return "config"
	}
	return vc.Directory
}

// IsFileExist
func (vc ViperConfig) IsFileExist() bool {
	return dir.IsExist(vc.getConfigFilePath())
}

// RemoveFile remove config file
func (vc ViperConfig) RemoveFile() error {
	return dir.Remove(vc.getConfigFilePath())
}

// RemoveDir remove config dir
func (vc ViperConfig) RemoveDir() error {
	return os.RemoveAll(filepath.Dir(vc.getConfigFilePath()))
}

// Recover recover config file content
func (vc ViperConfig) Recover(b []byte) error {
	_, err := dir.WriteBytes(vc.getConfigFilePath(), b)
	return err
}

// Init
func Init(vc ViperConfig) error {
	if vc.Name == "" {
		return ErrEmptyName
	}

	if vc.Type == "" {
		vc.Type = "yaml"
	}

	vc.Directory = vc.getConfigFileDir()

	filePath := vc.getConfigFilePath()

	vi := viper.New()
	vi.SetConfigName(vc.Name)
	vi.SetConfigType(vc.Type)
	vi.AddConfigPath(vc.Directory)

	isExist := dir.IsExist(filePath)
	if !isExist {
		if vc.Directory != "./" {
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

	if err := vc.Watch(vi); err != nil {
		return fmt.Errorf("watch config fail: %w ", err)
	}

	return nil
}
