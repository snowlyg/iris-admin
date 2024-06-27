package viper_server

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/spf13/viper"
)

var (
	ErrEmptyName = errors.New("config'name can't be empty value")
)

type ViperConfig struct {
	Debug     bool
	Directory string
	AbPath    string
	Name      string
	Type      string
	Default   []byte
	Watch     func(*viper.Viper) error
}

// getConfigFilePath
func (vc ViperConfig) getConfigFilePath() string {
	if vc.AbPath == "" {
		vc.AbPath = dir.GetCurrentAbPath()
	}
	return filepath.Join(vc.AbPath, vc.Directory, str.Join(vc.Name, ".", vc.Type))
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

// Remove remove config file
func (vc ViperConfig) Remove() error {
	return dir.Remove(vc.getConfigFilePath())
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
	if vc.Debug {
		fmt.Printf("\nthis config file's path is [%s]\n", filePath)
	}

	vi := viper.New()
	if vc.Debug {
		fmt.Printf("this config file's type is [%s]\n", vc.Type)
	}
	vi.SetConfigName(vc.Name)
	vi.SetConfigType(vc.Type)
	vi.AddConfigPath(vc.Directory)

	isExist := dir.IsExist(filePath)
	if !isExist {
		if vc.Debug {
			fmt.Printf("this config [%s] is not exist\n", filePath)
		}
		if vc.Directory != "./" {
			err := dir.InsureDir(filepath.Dir(filePath))
			if err != nil {
				return fmt.Errorf("create dir %s fail : %v", filePath, err)
			}
		}

		// ReadConfig
		if err := vi.ReadConfig(bytes.NewBuffer(vc.Default)); err != nil {
			if vc.Debug {
				fmt.Println(string(vc.Default))
			}
			return fmt.Errorf("read default config fail : %w ", err)
		}

		// WriteConfigAs
		if err := vi.WriteConfigAs(filePath); err != nil {
			return fmt.Errorf("write config to path fail: %w ", err)
		}

	} else {
		if vc.Debug {
			fmt.Printf("this config file [%s] is existed\n", filePath)
		}
		vi.SetConfigFile(filePath)
		err := vi.ReadInConfig()
		if err != nil {
			return fmt.Errorf("read config fail: %w ", err)
		}
	}

	err := vc.Watch(vi)
	if err != nil {
		return err
	}

	return nil
}
