package viper_server

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/spf13/viper"
)

var (
	ErrEmptyName = errors.New("配置文件名称不能为空")
)

type ViperConfig struct {
	Directory string
	Name      string
	Type      string
	Default   []byte
	Watch     func(*viper.Viper) error
}

// getConfigFilePath 获取配置文件路径
func (vc ViperConfig) getConfigFilePath() string {
	return str.Join(vc.Directory, "/", vc.Name, ".", vc.Type)
}

// getConfigFileDir 获取配置文件目录
func (vc ViperConfig) GetConfigFileDir() string {
	if vc.Directory == "" {
		return "config"
	}
	return vc.Directory
}

// Init 初始化系统配置
// - 第一次初始化系统配置，会自动生成配置文件 config.yaml 和 casbin 的规则文件 rbac_model.conf
// - 热监控系统配置项，如果发生变化会重写配置文件内的配置项
func Init(viperConfig ViperConfig) error {
	if viperConfig.Name == "" {
		return ErrEmptyName
	}

	if viperConfig.Type == "" {
		viperConfig.Type = "yaml"
	}

	viperConfig.Directory = viperConfig.GetConfigFileDir()

	fileName := viperConfig.getConfigFilePath()
	fmt.Printf("您的配置文件路径为 [%s]\n", fileName)

	vi := viper.New()
	fmt.Printf("您的配置文件类型为 [%s]\n", viperConfig.Type)
	vi.SetConfigType(viperConfig.Type)
	vi.AddConfigPath(viperConfig.Directory)

	isExist := dir.IsExist(fileName)
	if !isExist { //没有配置文件，写入默认配置
		fmt.Printf("您的配置文件 [%s] 不存在\n", fileName)
		if viperConfig.Directory != "./" {
			err := dir.InsureDir(viperConfig.Directory)
			if err != nil {
				return fmt.Errorf("新建 %s 目录失败 %v", fileName, err)
			}
		}

		// ReadConfig 读取配置文件到 vi 中
		if err := vi.ReadConfig(bytes.NewBuffer(viperConfig.Default)); err != nil {
			return fmt.Errorf("读取默认配置文件错误: %w ", err)
		}

		// WriteConfigAs 将配置写入本地文件中
		if err := vi.WriteConfigAs(fileName); err != nil {
			return fmt.Errorf("写入配置文件错误: %w ", err)
		}

	} else {
		fmt.Printf("您的配置文件 [%s] 已存在\n", fileName)
		// 存在配置文件，读取配置文件内容
		vi.SetConfigFile(fileName)
		err := vi.ReadInConfig()
		if err != nil {
			return fmt.Errorf("读取配置错误: %w ", err)
		}
	}

	err := viperConfig.Watch(vi)
	if err != nil {
		return err
	}

	return nil
	// 	casbinPath := filepath.Join(dir.GetCurrentAbPath(), g.CasbinFileName)
	// 	fmt.Printf("casbin rbac_model.conf 位于： %s\n\n", casbinPath)
	// 	if !dir.IsExist(casbinPath) { // casbin rbac_model.conf 文件
	// 		var rbacModelConf = []byte(`[request_definition]
	// r = sub, obj, act

	// [policy_definition]
	// p = sub, obj, act

	// [role_definition]
	// g = _, _

	// [policy_effect]
	// e = some(where (p.eft == allow))

	// [matchers]
	// m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")`)
	// 		_, err = dir.WriteBytes(casbinPath, rbacModelConf)
	// 		if err != nil {
	// 			panic(fmt.Errorf("初始化 casbin rbac_model.conf 文件错误: %w ", err))
	// 		}
	// 	}
}
