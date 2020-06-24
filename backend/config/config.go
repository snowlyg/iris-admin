package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

var Config = struct {
	HTTPS      bool   `default:"false" env:"HTTPS"`
	Certpath   string `default:"" env:"Certpath"`
	Certkey    string `default:"" env:"Certkey"`
	Port       uint   `default:"80" env:"PORT"`
	Host       string `default:"" env:"Host"`
	RecordPath string `default:"" env:"RecordPath"`
	Admin      struct {
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
	configPath := filepath.Join(Root, "config/application.yml")
	if len(os.Getenv("GOPATH")) == 0 {
		configPath = "config/application.yml"
	}

	if err := configor.Load(&Config, configPath); err != nil {
		panic(err)
	}

}
