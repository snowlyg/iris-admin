package conf

import "fmt"

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
