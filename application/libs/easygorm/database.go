package easygorm

import (
	"errors"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var easyGorm *DBServer

// DBServer 简单便捷的使用 gorm
type DBServer struct {
	conn     string
	db       *gorm.DB
	enforcer *casbin.Enforcer
	*Config
}

type Casbin struct {
	Path   string
	Prefix string
}

// Config 设置属性
type Config struct {
	GormConfig *gorm.Config
	Adapter    string        // 类型
	Conn       string        // 名称
	Models     []interface{} // 模型数据
	Casbin     *Casbin
}

func Init(c *Config) error {
	if easyGorm != nil {
		return nil
	}

	easyGorm = &DBServer{
		Config: &Config{
			GormConfig: &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   "g_",  // 表名前缀，`User` 的表名应该是 `t_users`
					SingularTable: false, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
				},
			},
			Adapter: "mysql", // 类型
			Conn:    "",
			Casbin: &Casbin{
				Path:   "",
				Prefix: "casbin",
			},
			Models: nil,
		},
	}

	if c != nil {
		easyGorm.Config = c
	}

	err := easyGorm.initGormDb()
	if err != nil {
		return errors.New(fmt.Sprintf("getGormDb err: : %+v", err))
	}
	if easyGorm.db == nil {
		return errors.New("数据库初始化失败")
	}

	if len(easyGorm.Config.Models) > 0 {
		err = Migrate(easyGorm.Config.Models)
		if err != nil {
			return errors.New(fmt.Sprintf("AutoMigrate err: : %+v", err))
		}
	}

	// 没有 CasbinPath 不使用 casbin
	if easyGorm.Casbin == nil {
		return errors.New("casbin 设置不能为空")
	}

	err = easyGorm.initEnforcer()
	if err != nil {
		return errors.New(fmt.Sprintf("getEnforcer err: : %+v", err))
	}

	if easyGorm.enforcer == nil {
		return errors.New("casbin 初始化失败")
	}

	return nil
}

// Migrate 迁移数据表
func Migrate(models []interface{}) error {
	err := easyGorm.db.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}

// initGormDb
func (db *DBServer) initGormDb() error {
	var err error
	var dialector gorm.Dialector
	if db.Config.Adapter == "mysql" {
		dialector = mysql.Open(db.Config.Conn)
	} else if db.Config.Adapter == "postgres" {
		dialector = postgres.Open(db.Config.Conn)
	} else if db.Config.Adapter == "sqlite3" {
		dialector = sqlite.Open(db.Config.Conn)
	} else {
		return errors.New("not supported database adapter")
	}

	//&gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		TablePrefix:   db.Config.Prefix, // 表名前缀，`User` 的表名应该是 `t_users`
	//		SingularTable: false,            // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
	//	},
	//}
	db.db, err = gorm.Open(dialector, db.Config.GormConfig)
	if err != nil {
		return err
	}

	err = db.db.Use(
		dbresolver.Register(
			dbresolver.Config{ /* xxx */ }).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)
	if err != nil {
		return err
	}

	db.db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: false})
	return nil
}

// initEnforcer
func (db *DBServer) initEnforcer() error {
	c, err := gormadapter.NewAdapterByDBUseTableName(db.db, db.Config.Casbin.Prefix, "casbin_rule") // Your driver and data source.
	if err != nil {
		return err
	}

	db.enforcer, err = casbin.NewEnforcer(db.Config.Casbin.Path, c)
	if err != nil {
		return err
	}

	err = db.enforcer.LoadPolicy()
	if err != nil {
		return err
	}

	return nil
}

func GetEasyGormDb() *gorm.DB {
	return easyGorm.db
}

func GetEasyGormEnforcer() *casbin.Enforcer {
	return easyGorm.enforcer
}
