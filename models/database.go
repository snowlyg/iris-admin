package models

import (
	"fmt"
	"os"
	"strings"

	"IrisAdminApi/transformer"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var Db *gorm.DB
var Enforcer *casbin.Enforcer
var err error

/**
*设置数据库连接
*@param diver string
 */
func Register(rc *transformer.Conf) {
	var c *gormadapter.Adapter

	if isTestEnv() {
		c, err = gormadapter.NewAdapter(rc.Sqlite.DirverName, rc.Sqlite.Connect) // Your driver and data source.
		if err != nil {
			panic(fmt.Sprintf("NewAdapter 错误: %v", err))
		}

		Db, err = gorm.Open(rc.Sqlite.DirverName, rc.Sqlite.Connect)
		if err != nil {
			panic(fmt.Sprintf("gorm open 错误: %v", err))
		}
	} else {
		c, err = gormadapter.NewAdapter(rc.Mysql.DirverName, rc.Mysql.Connect) // Your driver and data source.
		if err != nil {
			panic(fmt.Sprintf("NewAdapter 错误: %v", err))
		}

		Db, err = gorm.Open(rc.Mysql.DirverName, rc.Mysql.Connect+"/"+rc.Mysql.Name+"?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			panic(fmt.Sprintf("gorm open 错误: %v", err))
		}
	}

	Enforcer, err = casbin.NewEnforcer("./config/rbac_model.conf", c)
	if err != nil {
		panic(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}
	_ = Enforcer.LoadPolicy()

}

//获取程序运行环境
// 根据程序运行路径后缀判断
//如果是 test 就是测试环境
func isTestEnv() bool {
	files := os.Args
	for _, v := range files {
		if strings.Contains(v, "test") {
			return true
		}
	}
	return false
}

/**
 * 获取列表
 * @method MGetAll
 * @param  {[type]} string string    [description]
 * @param  {[type]} orderBy string    [description]
 * @param  {[type]} relation string    [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAll(string, orderBy string, offset, limit int) *gorm.DB {

	if len(orderBy) > 0 {
		Db.Order(orderBy + "desc")
	} else {
		Db.Order("created_at desc")
	}

	if len(string) > 0 {
		Db.Where("name LIKE  ?", "%"+string+"%")
	}

	if offset > 0 {
		Db.Offset((offset - 1) * limit)
	}

	if limit > 0 {
		Db.Limit(limit)
	}

	return Db
}
