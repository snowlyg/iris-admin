package easygorm

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"os"
	"path/filepath"
	"time"

	_ "gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Egm *EasyGorm

// EasyGorm 简单便捷的使用 gorm
type EasyGorm struct {
	Conn     string
	Db       *gorm.DB
	Enforcer *casbin.Enforcer
	*Config
}

// Config 设置属性
type Config struct {
	GormConfig      *gorm.Config
	Adapter         string // 类型
	Name            string // 名称
	Username        string // 用户名
	Pwd             string // 密码
	Host            string // 地址
	Port            int64  // 端口
	CasbinModelPath string // casbin 模型规则路径
	SqlitePath      string // sqlite 模型路径
	Debug           bool   // 调试
}

func Init(c *Config) {
	Egm = new(EasyGorm)
	Egm.Config = &Config{
		GormConfig: &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "sg_", // 表名前缀，`User` 的表名应该是 `t_users`
				SingularTable: false, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
			},
		},
		Adapter:         "mysql", // 类型
		Name:            "",
		Username:        "root",      // 用户名
		Pwd:             "",          // 密码
		Host:            "127.0.0.1", // 地址
		Port:            3306,        // 端口
		CasbinModelPath: "",          // casbin 模型规则路径
	}

	if c != nil {
		Egm.Config = c
	}

	if Egm.Config != nil {
		Egm.getGormDb()
		if Egm.CasbinModelPath != "" {
			Egm.getEnforcer()
		}

		if Egm.Config.Debug {
			fmt.Println(fmt.Sprintf("easygorm config: : %+v", Egm.Config))
		}
	}
}

// getConn 获取链接
func (db *EasyGorm) getConn() string {
	if db.Config.Adapter == "mysql" {
		return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", db.Config.Username, db.Config.Pwd, db.Config.Host, db.Config.Port, db.Config.Name)
	} else if db.Config.Adapter == "postgres" {
		return fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", db.Config.Username, db.Config.Pwd, db.Config.Host, db.Config.Name)
	} else if db.Config.Adapter == "sqlite3" {
		return filepath.Join(cwd(), db.Config.Name) + ".db"
	} else {
		fmt.Println(errors.New("not supported database adapter"))
	}

	return ""
}

// getGormDb
func (db *EasyGorm) getGormDb() {
	var err error
	var dialector gorm.Dialector
	if db.Config.Adapter == "mysql" {
		dialector = mysql.Open(db.getConn())
	} else if db.Config.Adapter == "postgres" {
		dialector = postgres.Open(db.getConn())
	} else if db.Config.Adapter == "sqlite3" {
		dialector = sqlite.Open(db.getConn())
	} else {
		fmt.Println(errors.New("not supported database adapter"))
	}

	if db.Config.Debug {
		fmt.Println(fmt.Sprintf("Conn: : %s", db.getConn()))
	}

	//&gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		TablePrefix:   db.Config.Prefix, // 表名前缀，`User` 的表名应该是 `t_users`
	//		SingularTable: false,            // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
	//	},
	//}
	db.Db, err = gorm.Open(dialector, db.Config.GormConfig)
	if err != nil {
		fmt.Println(err)
	}

	_ = db.Db.Use(
		dbresolver.Register(
			dbresolver.Config{ /* xxx */ }).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)
	db.Db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: false})

}

// getEnforcer
func (db *EasyGorm) getEnforcer() {
	var err error

	c, err := gormadapter.NewAdapter(db.Config.Adapter, db.getConn(), true) // Your driver and data source.
	if err != nil {
		fmt.Println(fmt.Sprintf("NewAdapter 错误: %v,Path: %s", err, db.getConn()))
	}

	db.Enforcer, err = casbin.NewEnforcer(db.Config.CasbinModelPath, c)
	if err != nil {
		fmt.Println(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}

	_ = db.Enforcer.LoadPolicy()
}

// cwd 获取项目路径
func cwd() string {
	// 兼容 travis 集成测试
	if os.Getenv("TRAVIS_BUILD_DIR") != "" {
		return os.Getenv("TRAVIS_BUILD_DIR")
	}

	path, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(path)
}
