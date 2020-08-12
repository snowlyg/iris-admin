package sysinit

import (
	"errors"
	"fmt"
	"github.com/snowlyg/IrisAdminApi/server/libs"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/snowlyg/IrisAdminApi/server/config"
)

var Enforcer *casbin.Enforcer

func init() {

	var err error
	var conn string
	if config.Config.DB.Adapter == "mysql" {
		conn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", config.Config.DB.User, config.Config.DB.Password, config.Config.DB.Host, config.Config.DB.Port, config.Config.DB.Name)
	} else if config.Config.DB.Adapter == "postgres" {
		conn = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", config.Config.DB.User, config.Config.DB.Password, config.Config.DB.Host, config.Config.DB.Name)
	} else if config.Config.DB.Adapter == "sqlite3" {
		conn = libs.DBFile()
	} else {
		panic(errors.New("not supported database adapter"))
	}

	if len(conn) == 0 {
		panic(fmt.Sprintf("数据链接不可用: %s", conn))
	}

	c, err := gormadapter.NewAdapter(config.Config.DB.Adapter, conn, true) // Your driver and data source.
	if err != nil {
		panic(fmt.Sprintf("NewAdapter 错误: %v", err))
	}

	casbinModelPath := filepath.Join(libs.CWD(), "rbac_model.conf")
	fmt.Println(casbinModelPath)
	Enforcer, err = casbin.NewEnforcer(casbinModelPath, c)
	if err != nil {
		panic(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}

	_ = Enforcer.LoadPolicy()

}
