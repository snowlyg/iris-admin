package database

import (
	"encoding/json"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG = Mysql{
	Path:         "127.0.0.1:3306",
	Config:       "charset=utf8mb4&parseTime=True&loc=Local",
	DbName:       "iris-admin",
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
	DbName       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`
	LogZap       string `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"` //silent,error,warn,info,zap
}

// Dsn return mysql dsn
func (m *Mysql) Dsn() string {
	return fmt.Sprintf("%s%s?%s", m.BaseDsn(), m.DbName, m.Config)
}

// Dsn return
func (m *Mysql) BaseDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/", m.Username, m.Password, m.Path)
}

// IsExist config file is exist
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove remove config file
func Remove() error {
	return getViperConfig().Remove()
}

// Recover
func Recover() error {
	b, err := json.Marshal(CONFIG)
	if err != nil {
		return err
	}
	return getViperConfig().Recover(b)
}

// getViperConfig get viper config
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
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			// watch config file change
			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("Config file changed:", e.Name)
			})
			vi.WatchConfig()
			return nil
		},
		//
		Default: []byte(`
{
	"path": "` + CONFIG.Path + `",
	"config": "` + CONFIG.Config + `",
	"db-name": "` + CONFIG.DbName + `",
	"username": "` + CONFIG.Username + `",
	"password": "` + CONFIG.Password + `",
	"max-idle-conns": ` + mxIdleConns + `,
	"max-open-conns": ` + mxOpenConns + `,
	"log-mode": ` + logMode + `,
	"log-zap": "` + CONFIG.LogZap + `"
}`),
	}
}
