package models

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"IrisAdminApi/transformer"
	"IrisAdminApi/validates"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var Db *gorm.DB
var Enforcer *casbin.Enforcer
var err error
var c *gormadapter.Adapter
var dirverName string
var casbinConn string
var conn string

/**
*设置数据库连接
*@param diver string
 */
func Register(rc *transformer.Conf) {
	if rc.App.DirverType == "Sqlite" {
		dirverName = rc.Sqlite.DirverName
		if isTestEnv() {
			casbinConn = rc.Sqlite.TConnect
			conn = rc.Sqlite.TConnect
		} else {
			casbinConn = rc.Sqlite.Connect
			conn = rc.Sqlite.Connect
		}

	} else if rc.App.DirverType == "Mysql" {
		dirverName = rc.Mysql.DirverName
		if isTestEnv() {
			casbinConn = rc.Mysql.Connect
			conn = rc.Mysql.Connect + rc.Mysql.TName + "?charset=utf8&parseTime=True&loc=Local"
		} else {
			casbinConn = rc.Mysql.Connect
			conn = rc.Mysql.Connect + rc.Mysql.Name + "?charset=utf8&parseTime=True&loc=Local"
		}
	}

	Db, err = gorm.Open(dirverName, conn)
	if err != nil {
		color.Red(fmt.Sprintf("gorm open 错误: %v", err))
	}

	c, err = gormadapter.NewAdapter(dirverName, casbinConn) // Your driver and data source.
	if err != nil {
		color.Red(fmt.Sprintf("NewAdapter 错误: %v", err))
	}

	Enforcer, err = casbin.NewEnforcer("./config/rbac_model.conf", c)
	if err != nil {
		color.Red(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}
	_ = Enforcer.LoadPolicy()

}

/**
*初始化系统 账号 权限 角色
 */
func CreateSystemData(rc *transformer.Conf, perms []*validates.PermissionRequest) {
	if rc.App.CreateSysData {
		permIds := CreateSystemAdminPermission(perms) //初始化权限
		role := CreateSystemAdminRole(permIds)        //初始化角色
		if role.ID != 0 {
			CreateSystemAdmin(role.ID, rc) //初始化管理员
		}
	}
}

func IsNotFound(err error) {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); !ok && err != nil {
		color.Red(fmt.Sprintf("error :%v \n ", err))
	}
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
