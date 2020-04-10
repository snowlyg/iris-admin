package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

var Config = struct {
	HTTPS    bool   `default:"false" env:"HTTPS"`
	Certpath string `default:"" env:"Certpath"`
	Certkey  string `default:"" env:"Certkey"`
	Port     uint   `default:"8085" env:"PORT"`
	Host     string `default:"127.0.0.1" env:"Host"`
	Admin    struct {
		UserName        string `env:"AdminUserName" default:"username"`
		Name            string `env:"AdminName" default:"name"`
		Pwd             string `env:"AdminPwd" default:"123456"`
		RoleName        string `env:"AdminRoleName" default:"superadmin_role"`
		RoleDisplayName string `env:"TenantRoleDisplayName" default:"商户管理员"`
	}
	DB struct {
		Name     string `env:"DBName" default:"goirisadminapi"`
		Adapter  string `env:"DBAdapter" default:"mysql"`
		Host     string `env:"DBHost" default:"mysql"`
		Port     string `env:"DBPort" default:"3306"`
		User     string `env:"DBUser" default:"root"`
		Password string `env:"DBPassword" default:"123456"`
	}
}{}

var Root = os.Getenv("GOPATH") + "/src/github.com/snowlyg/IrisAdminApi/backend/"

func init() {
	if err := configor.Load(&Config, filepath.Join(Root, "config/application.yml")); err != nil {
		panic(err)
	}
}
