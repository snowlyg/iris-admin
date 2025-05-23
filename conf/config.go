package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	ConfigType = "json"   // config's type
	ConfigDir  = "config" // config's dir
)

var (
	mysqlAddrKey = "IRIS_ADMIN_MYSQL_ADDR"
	mysqlPwdKey  = "IRIS_ADMIN_MYSQL_PWD"
	mysqlNameKey = "IRIS_ADMIN_MYSQL_NAME"
	webAddrKey   = "IRIS_ADMIN_WEB_ADDR"
)

func NewConf() *Conf {
	c := &Conf{
		Locale:         "zh",
		FileMaxSize:    1024,   // upload file size limit 1024M
		SessionTimeout: 172800, // session timeout after 4 months
		CorsConf: CorsConf{
			AccessOrigin:        "*",
			AccessHeaders:       "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id",
			AccessMethods:       "POST,GET,OPTIONS,DELETE,PUT",
			AccessExposeHeaders: "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type",
			AccessCredentials:   "true",
		},
		Except: Route{
			Uri:    "",
			Method: "",
		},
		System: System{
			Tls:        false,
			GinMode:    gin.ReleaseMode,
			Level:      "debug",
			Addr:       "127.0.0.1:8080",
			TimeFormat: "2006-01-02 15:04:05",
		},
		Limit: Limit{
			Disable: true,
			Limit:   0,
			Burst:   5,
		},
		Captcha: Captcha{
			KeyLong:   0,
			ImgWidth:  240,
			ImgHeight: 80,
		},
		Mysql: &Mysql{
			Path:         "127.0.0.1:3306",
			Config:       "charset=utf8mb4&parseTime=True&loc=Local",
			DbName:       "iris-admin",
			Username:     "root",
			Password:     "",
			MaxIdleConns: 0,
			MaxOpenConns: 0,
			LogMode:      false,
			LogZap:       "error",
		},
		Operate: Operate{
			Except: Route{
				Uri:    "api/v1/upload;api/v1/upload",
				Method: "post;put",
			},
			Include: Route{
				Uri:    "api/v1/menus",
				Method: "get",
			},
		},
	}
	mysqlAddr := strings.TrimSpace(os.Getenv(mysqlAddrKey))
	mysqlPwd := strings.TrimSpace(os.Getenv(mysqlPwdKey))
	mysqlName := strings.TrimSpace(os.Getenv(mysqlNameKey))
	webAddr := strings.TrimSpace(os.Getenv(webAddrKey))
	if mysqlAddr != "" {
		c.Mysql.Path = mysqlAddr
	}
	if mysqlPwd != "" {
		c.Mysql.Password = mysqlPwd
	}
	if mysqlName != "" {
		c.Mysql.Username = mysqlName
	}
	if webAddr != "" {
		c.System.Addr = webAddr
	}
	if c.Mysql.Path == "" || c.Mysql.Password == "" || c.Mysql.DbName == "" {
		log.Printf("mysql driver config empty,you can set env %s %s %s to change it.\n", mysqlAddrKey, mysqlPwdKey, mysqlNameKey)
	}
	return c
}

type Conf struct {
	Locale         string   `mapstructure:"locale" json:"locale" yaml:"locale"`
	FileMaxSize    int64    `mapstructure:"file-max-size" json:"file-max-size" yaml:"file-max-siz"`
	SessionTimeout int64    `mapstructure:"session-timeout" json:"session-timeout" yaml:"session-timeout"`
	Except         Route    `mapstructure:"except" json:"except" yaml:"except"`
	System         System   `mapstructure:"system" json:"system" yaml:"system"`
	Limit          Limit    `mapstructure:"limit" json:"limit" yaml:"limit"`
	Captcha        Captcha  `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	CorsConf       CorsConf `mapstructure:"cors" json:"cors" yaml:"cors"`
	Mysql          *Mysql   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Operate        Operate  `mapstructure:"operate" json:"operate" yaml:"operate"`
}

type Route struct {
	Uri    string `mapstructure:"uri" json:"uri" yaml:"uri"`
	Method string `mapstructure:"method" json:"method" yaml:"method"`
}

type Captcha struct {
	KeyLong   int `mapstructure:"key-long" json:"key-long" yaml:"key-long"`
	ImgWidth  int `mapstructure:"img-width" json:"img-width" yaml:"img-width"`
	ImgHeight int `mapstructure:"img-height" json:"img-height" yaml:"img-height"`
}

type Limit struct {
	Disable bool    `mapstructure:"disable" json:"disable" yaml:"disable"`
	Limit   float64 `mapstructure:"limit" json:"limit" yaml:"limit"`
	Burst   int     `mapstructure:"burst" json:"burst" yaml:"burst"`
}

type System struct {
	GinMode    string `mapstructure:"gin-mode" json:"gin-mode" yaml:"gin-mode"`
	Tls        bool   `mapstructure:"tls" json:"tls" yaml:"tls"`
	Level      string `mapstructure:"level" json:"level" yaml:"level"` // debug,release,test
	Addr       string `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType     string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`
	TimeFormat string `mapstructure:"time-format" json:"time-format" yaml:"time-format"`
}

// SetDefaultAddrAndTimeFormat
func (conf *Conf) SetDefaultAddrAndTimeFormat() {
	if conf.System.Addr == "" {
		conf.System.Addr = "127.0.0.1:8080"
	}

	if conf.System.TimeFormat == "" {
		conf.System.TimeFormat = "2006-01-02 15:04:05"
	}
}

// // toStaticUrl
// func (conf *Conf) toStaticUrl(uri string) string {
// 	path := filepath.Join(conf.System.Addr, conf.System.StaticPrefix, uri)
// 	if conf.System.Tls {
// 		return filepath.ToSlash(str.Join("https://", path))
// 	}
// 	return filepath.ToSlash(str.Join("http://", path))
// }

// IsExist config file is exist
func (conf *Conf) IsExist() bool {
	return conf.getViperConfig().IsExist()
}

// RemoveFile remove config file
func (conf *Conf) RemoveFile() error {
	return conf.getViperConfig().RemoveFile()
}

// Recover
func (conf *Conf) Recover() error {
	conf.newRbacModel()
	b, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return fmt.Errorf("iris-admin recover config faild:%w", err)
	}
	return conf.getViperConfig().Recover(b)
}

// getViperConfig get viper config
func (conf *Conf) getViperConfig() *ViperConf {
	maxSize := strconv.FormatInt(conf.FileMaxSize, 10)
	sessionTimeout := strconv.FormatInt(conf.SessionTimeout, 10)
	keyLong := strconv.FormatInt(int64(conf.Captcha.KeyLong), 10)
	imgWidth := strconv.FormatInt(int64(conf.Captcha.ImgWidth), 10)
	imgHeight := strconv.FormatInt(int64(conf.Captcha.ImgHeight), 10)
	limit := strconv.FormatInt(int64(conf.Limit.Limit), 10)
	burst := strconv.FormatInt(int64(conf.Limit.Burst), 10)
	disable := strconv.FormatBool(conf.Limit.Disable)
	tls := strconv.FormatBool(conf.System.Tls)

	mxIdleConns := fmt.Sprintf("%d", conf.Mysql.MaxIdleConns)
	mxOpenConns := fmt.Sprintf("%d", conf.Mysql.MaxOpenConns)
	logMode := fmt.Sprintf("%t", conf.Mysql.LogMode)

	configName := "iris_admin"
	return &ViperConf{
		dir:  ConfigDir,
		name: configName,
		t:    ConfigType,
		watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&conf); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			// watch config file change
			vi.SetConfigName(configName)
			return nil
		},
		//
		Default: []byte(`
{
	"locale": "` + conf.Locale + `",
	"file-max-size": ` + maxSize + `,
	"session-timeout": ` + sessionTimeout + `,
	"except":
		{ 
			"uri": "` + conf.Except.Uri + `",
			"method": "` + conf.Except.Method + `"
		},
	"cors":
		{ 
			"access-origin": "` + conf.CorsConf.AccessOrigin + `",
			"access-headers": "` + conf.CorsConf.AccessHeaders + `",
			"access-methods": "` + conf.CorsConf.AccessMethods + `",
			"access-expose-headers": "` + conf.CorsConf.AccessExposeHeaders + `",
			"access-credentials": "` + conf.CorsConf.AccessCredentials + `"
		},
	"captcha":
		{
		"key-long": ` + keyLong + `,
		"img-width": ` + imgWidth + `,
		"img-height": ` + imgHeight + `
		},
	"limit":
		{
			"limit": ` + limit + `,
			"disable": ` + disable + `,
			"burst": ` + burst + `
		},
	"system":
		{
			"tls": ` + tls + `,
			"level": "` + conf.System.Level + `",
			"gin-mode": "` + conf.System.GinMode + `",
			"addr": "` + conf.System.Addr + `",
			"time-format": "` + conf.System.TimeFormat + `"
		},
		"mysql":
		{
			"path": "` + conf.Mysql.Path + `",
			"config": "` + conf.Mysql.Config + `",
			"db-name": "` + conf.Mysql.DbName + `",
			"username": "` + conf.Mysql.Username + `",
			"password": "` + conf.Mysql.Password + `",
			"max-idle-conns": ` + mxIdleConns + `,
			"max-open-conns": ` + mxOpenConns + `,
			"log-mode": ` + logMode + `,
			"log-zap": "` + conf.Mysql.LogZap + `"
		},
		"operate":
		{
			"except":{ 
				"uri": "` + conf.Operate.Except.Uri + `",
				"method": "` + conf.Operate.Except.Method + `"
			},	
			"include":
			{
				"uri": "` + conf.Operate.Include.Uri + `",
				"method": "` + conf.Operate.Include.Method + `"
			} 
		}
 }`),
	}
}
