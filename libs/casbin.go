package libs

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {

	var err error
	var conn string
	if Config.DB.Adapter == "mysql" {
		conn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", Config.DB.User, Config.DB.Password, Config.DB.Host, Config.DB.Port, Config.DB.Name)
	} else if Config.DB.Adapter == "postgres" {
		conn = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", Config.DB.User, Config.DB.Password, Config.DB.Host, Config.DB.Name)
	} else if Config.DB.Adapter == "sqlite3" {
		conn = DBFile()
	} else {
		logger.Println(errors.New("not supported database adapter"))
	}

	if len(conn) == 0 {
		logger.Println(fmt.Sprintf("数据链接不可用: %s", conn))
	}

	c, err := gormadapter.NewAdapter(Config.DB.Adapter, conn, true) // Your driver and data source.
	if err != nil {
		logger.Println(fmt.Sprintf("NewAdapter 错误: %v,Path: %s", err, conn))
	}

	casbinModelPath := filepath.Join(CWD(), "rbac_model.conf")
	Enforcer, err = casbin.NewEnforcer(casbinModelPath, c)
	if err != nil {
		logger.Println(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}

	_ = Enforcer.LoadPolicy()

}
