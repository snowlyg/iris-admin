package sysinit

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/snowlyg/IrisAdminApi/server/config"
	"github.com/snowlyg/IrisAdminApi/server/libs"
)

var (
	Db *gorm.DB
)

func init() {

	var err error
	var conn string
	if config.Config.DB.Adapter == "mysql" {
		conn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", config.Config.DB.User, config.Config.DB.Password, config.Config.DB.Host, config.Config.DB.Port, config.Config.DB.Name)
	} else if config.Config.DB.Adapter == "postgres" {
		conn = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", config.Config.DB.User, config.Config.DB.Password, config.Config.DB.Host, config.Config.DB.Name)
	} else if config.Config.DB.Adapter == "sqlite3" {
		fmt.Println(libs.DBFile())
		conn = libs.DBFile()
	} else {
		panic(errors.New("not supported database adapter"))
	}

	Db, err = gorm.Open(config.Config.DB.Adapter, conn)
	if err != nil {
		panic(err)
	}

	gorm.DefaultTableNameHandler = func(Db *gorm.DB, defaultTableName string) string {
		return "iris_" + defaultTableName
	}

	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)
}
