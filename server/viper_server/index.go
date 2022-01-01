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
	ErrEmptyName = errors.New("配置文件名称不能为空")
)

type ViperConfig struct {
	Debug     bool
	Directory string
	Name      string
	Type      string
	Default   []byte
	Watch     func(*viper.Viper) error
}

// getConfigFilePath 获取配置文件路径
func (vc ViperConfig) getConfigFilePath() string {
	return filepath.Join(dir.GetCurrentAbPath(), vc.Directory, str.Join(vc.Name, ".", vc.Type))
}

// getConfigFileDir 获取配置文件目录
func (vc ViperConfig) getConfigFileDir() string {
	if vc.Directory == "" {
		return "config"
	}
	return vc.Directory
}

// IsFileExist 获取配置文件目录
func (vc ViperConfig) IsFileExist() bool {
	return dir.IsExist(vc.getConfigFilePath())
}

// Remove 删除配置文件
func (vc ViperConfig) Remove() error {
	return dir.Remove(vc.getConfigFilePath())
}

// Init 初始化系统配置
// - 第一次初始化系统配置，会自动生成配置文件 config.yaml 和 casbin 的规则文件 rbac_model.conf
// - 热监控系统配置项，如果发生变化会重写配置文件内的配置项
func Init(vc ViperConfig) error {
	if vc.Name == "" {
		return ErrEmptyName
	}

	if vc.Type == "" {
		vc.Type = "yaml"
	}

	vc.Directory = vc.getConfigFileDir()

	fileName := vc.getConfigFilePath()
	if vc.Debug {
		fmt.Printf("您的配置文件路径为 [%s]\n", fileName)
	}

	vi := viper.New()
	if vc.Debug {
		fmt.Printf("您的配置文件类型为 [%s]\n", vc.Type)
	}
	vi.SetConfigType(vc.Type)
	vi.AddConfigPath(vc.Directory)

	isExist := dir.IsExist(fileName)
	if !isExist { //没有配置文件，写入默认配置
		if vc.Debug {
			fmt.Printf("您的配置文件 [%s/%s] 不存在\n", vc.Directory, fileName)
		}
		if vc.Directory != "./" {
			err := dir.InsureDir(filepath.Dir(fileName))
			if err != nil {
				return fmt.Errorf("新建 %s 目录失败 %v", fileName, err)
			}
		}

		// ReadConfig 读取配置文件到 vi 中
		if err := vi.ReadConfig(bytes.NewBuffer(vc.Default)); err != nil {
			if vc.Debug {
				fmt.Println(string(vc.Default))
			}
			return fmt.Errorf("读取默认配置文件错误: %w ", err)
		}

		// WriteConfigAs 将配置写入本地文件中
		if err := vi.WriteConfigAs(fileName); err != nil {
			return fmt.Errorf("写入配置文件错误: %w ", err)
		}

	} else {
		if vc.Debug {
			fmt.Printf("您的配置文件 [%s] 已存在\n", fileName)
		}
		// 存在配置文件，读取配置文件内容
		vi.SetConfigFile(fileName)
		err := vi.ReadInConfig()
		if err != nil {
			return fmt.Errorf("读取配置错误: %w ", err)
		}
	}

	err := vc.Watch(vi)
	if err != nil {
		return err
	}

	return nil
}
