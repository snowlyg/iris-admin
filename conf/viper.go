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
	dir     string
	name    string
	t       string
	Default []byte
	watch   func(*viper.Viper) error
}

// getConfPath
func (vc *ViperConf) getConfPath() string {
	if vc == nil {
		return ""
	}
	return filepath.Join(dir.GetCurrentAbPath(), vc.Dir(), str.Join(vc.name, ".", vc.t))
}

// Dir
func (vc *ViperConf) Dir() string {
	if vc.dir == "" {
		vc.dir = "config"
		return vc.dir
	}
	return vc.dir
}

// IsExist
func (vc *ViperConf) IsExist() bool {
	if vc == nil {
		return false
	}
	return dir.IsExist(vc.getConfPath())
}

// RemoveFile remove config file
func (vc *ViperConf) RemoveFile() error {
	if vc == nil {
		return e.ErrViperConfInvalid
	}
	d := filepath.Dir(vc.getConfPath())
	b := filepath.Base(d)
	if b != vc.Dir() {
		return nil
	}
	return dir.Remove(vc.getConfPath())
}

// RemoveDir remove config dir
func (vc *ViperConf) RemoveDir() error {
	if vc == nil {
		return e.ErrViperConfInvalid
	}
	d := filepath.Dir(vc.getConfPath())
	b := filepath.Base(d)
	if b != vc.Dir() {
		return fmt.Errorf("%s viper conf base '%s' want but get '%s'", d, b, vc.Dir())
	}
	return os.RemoveAll(d)
}

// Recover
func (vc *ViperConf) Recover(b []byte) error {
	if vc == nil {
		return e.ErrViperConfInvalid
	}
	_, err := dir.WriteBytes(vc.getConfPath(), b)
	return err
}

// NewViperConf
func NewViperConf(vc *ViperConf) error {
	if vc == nil {
		return e.ErrViperConfInvalid
	}
	if vc.name == "" {
		return e.ErrConfigNameEmpty
	}
	if vc.t == "" {
		vc.t = "yaml"
	}

	vc.dir = vc.Dir()
	filePath := vc.getConfPath()

	vi := viper.New()
	vi.SetConfigName(vc.name)
	vi.SetConfigType(vc.t)
	vi.AddConfigPath(vc.dir)
	isExist := dir.IsExist(filePath)
	if !isExist {
		if vc.Dir() != "./" {
			if err := dir.InsureDir(filepath.Dir(filePath)); err != nil {
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
