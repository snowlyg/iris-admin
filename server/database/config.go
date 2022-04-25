package database

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG = Mysql{
	Path:         "127.0.0.1:3306",
	Config:       "charset=utf8mb4&parseTime=True&loc=Local",
	Dbname:       "iris-admin",
	Username:     "root",
	Password:     "",
	MaxIdleConns: 0,
	MaxOpenConns: 0,
	LogMode:      false,
	LogZap:       "error",
}

type Mysql struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`
	LogZap       string `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"` //silent,error,warn,info,zap
}

// Dsn 获取 mysql dsn
func (m *Mysql) Dsn() string {
	return fmt.Sprintf("%s%s?%s", m.BaseDsn(), m.Dbname, m.Config)
}

// Dsn 获取 mysql dsn
func (m *Mysql) BaseDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/", m.Username, m.Password, m.Path)
}

// IsExist 配置文件是否存在
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove 删除配置文件
func Remove() error {
	err := getViperConfig().Remove()
	if err != nil {
		return err
	}
	return nil
}

// getViperConfig 获取初始化配置
func getViperConfig() viper_server.ViperConfig {
	configName := "mysql"
	mxIdleConns := fmt.Sprintf("%d", CONFIG.MaxIdleConns)
	mxOpenConns := fmt.Sprintf("%d", CONFIG.MaxOpenConns)
	logMode := fmt.Sprintf("%t", CONFIG.LogMode)
	return viper_server.ViperConfig{
		Debug:     true,
		Directory: g.ConfigDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("反序列化错误: %v", err)
			}
			// 监控配置文件变化
			vi.SetConfigName(configName)
			vi.WatchConfig()
			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("配置发生变化:", e.Name)
				if err := vi.Unmarshal(&CONFIG); err != nil {
					fmt.Printf("反序列化错误: %v \n", err)
				}
			})
			return nil
		},
		// 注意:设置默认配置值的时候,前面不能有空格等其他符号.必须紧贴左侧.
		Default: []byte(`
{
	"path": "` + CONFIG.Path + `",
	"config": "` + CONFIG.Config + `",
	"db-name": "` + CONFIG.Dbname + `",
	"username": "` + CONFIG.Username + `",
	"password": "` + CONFIG.Password + `",
	"max-idle-conns": ` + mxIdleConns + `,
	"max-open-conns": ` + mxOpenConns + `,
	"log-mode": ` + logMode + `,
	"log-zap": "` + CONFIG.LogZap + `"
}`),
	}
}
