package libs

import (
	"fmt"
	"github.com/snowlyg/blog/libs/logging"
	"path/filepath"
	"strings"

	"github.com/jinzhu/configor"
)

var Config = struct {
	LogLevel string `default:"info" env:"Loglevel"`
	HTTPS    bool   `default:"false" env:"Https"`
	Certpath string `default:"" env:"Certpath"`
	Certkey  string `default:"" env:"Certkey"`
	Port     int64  `default:"8085" env:"Port"`
	Host     string `default:"127.0.0.1" env:"Host"`
	Cache    struct {
		Driver string `env:"CacheDriver" default:"local"`
	}
	CasbinModel string `default:"" env:"CasbinModel"`
	Admin       struct {
		UserName        string `env:"AdminUserName" default:"username"`
		Name            string `env:"AdminName" default:"name"`
		Pwd             string `env:"AdminPwd" default:"123456"`
		RoleName        string `env:"AdminRoleName" default:"admin"`
		RoleDisplayName string `env:"TenantRoleDisplayName" default:"超级管理员"`
	}
	DB struct {
		Prefix   string `env:"DBPrefix" default:"iris_"`
		Name     string `env:"DBName" default:"blog"`
		Adapter  string `env:"DBAdapter" default:"mysql"`
		Host     string `env:"DBHost" default:"localhost"`
		Port     int64  `env:"DBPort" default:"3306"`
		User     string `env:"DBUser" default:"root"`
		Password string `env:"DBPassword" default:""`
	}
	Redis struct {
		Host string `env:"RedisHost" default:"localhost"`
		Port string `env:"RedisPort" default:"6379"`
		Pwd  string `env:"RedisPwd" default:""`
	}
	Limit struct {
		Disable bool    `env:"LimitDisable" default:"true"`
		Limit   float64 `env:"LimitLimit" default:"1"`
		Burst   int     `env:"LimitBurst" default:"5"`
	}
	Qiniu struct {
		Enable    bool   `env:"QiniuEnable" default:"false"`
		Host      string `env:"QiniuHost" default:""`
		Accesskey string `env:"QiniuAccesskey" default:""`
		Secretkey string `env:"QiniuSecretkey" default:""`
		Bucket    string `env:"QiniuBucket" default:""`
	}
}{}

func InitConfig(cp, cm string) {
	path := filepath.Join(CWD(), "application.yml")
	if cp != "" {
		path = cp
	}

	logging.Dbug.Debugf("配置YML文件路径：%v\n", path)
	if err := configor.Load(&Config, path); err != nil {
		logging.Err.Errorf("Config Path:%s ,Error:%+v\n", path, err)
		return
	}

	if cm != "" {
		Config.CasbinModel = cm
		return
	}

	if Config.CasbinModel == "" {
		Config.CasbinModel = filepath.Join(CWD(), "rbac_model.conf")
	}

	logging.Dbug.Debugf("配置项：%+v\n", Config)
}

func GetRedisUris() []string {
	addrs := make([]string, 0, 0)
	hosts := strings.Split(Config.Redis.Host, ";")
	ports := strings.Split(Config.Redis.Port, ";")
	for _, h := range hosts {
		for _, p := range ports {
			addrs = append(addrs, fmt.Sprintf("%s:%s", h, p))
		}
	}
	return addrs
}
