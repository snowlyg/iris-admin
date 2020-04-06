package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

var Config = struct {
	HTTPS  bool   `default:"false" env:"HTTPS"`
	Port   uint   `default:"7000" env:"PORT"`
	Host   string `default:"localhost" env:"Host"`
	Tenant struct {
		RoleName        string `env:"TenantRoleName" default:"tenant_role"`
		RoleDisplayName string `env:"TenantRoleDisplayName" default:"超级管理员"`
	}
	Admin struct {
		UserName        string `env:"AdminUserName" default:"username"`
		Name            string `env:"AdminName" default:"name"`
		Pwd             string `env:"AdminPwd" default:"123456"`
		RoleName        string `env:"AdminRoleName" default:"superadmin_role"`
		RoleDisplayName string `env:"TenantRoleDisplayName" default:"商户管理员"`
	}
	DB struct {
		Name     string `env:"DBName" default:"qor_example"`
		Adapter  string `env:"DBAdapter" default:"mysql"`
		Host     string `env:"DBHost" default:"mysql"`
		Port     string `env:"DBPort" default:"3306"`
		User     string `env:"DBUser" default:"root"`
		Password string `env:"DBPassword" default:""`
	}
}{}

var Root = os.Getenv("GOPATH") + "/src/github.com/snowlyg/IrisAdminApi"

func init() {
	if err := configor.Load(&Config, filepath.Join(Root, "config/application.yml")); err != nil {
		panic(err)
	}
}
