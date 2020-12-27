package easygorm

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"time"

	"gorm.io/gorm"
)

var EasyGorm *DBServer

// DBServer 简单便捷的使用 gorm
type DBServer struct {
	conn string
	*gorm.DB
	*casbin.Enforcer
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
	*Casbin
}

func Init(c *Config) error {
	if EasyGorm != nil {
		return nil
	}

	EasyGorm = &DBServer{
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
		EasyGorm.Config = c
	}

	err := EasyGorm.getGormDb()
	if err != nil {
		return errors.New(fmt.Sprintf("getGormDb err: : %+v", err))
	}
	if EasyGorm.DB == nil {
		return errors.New("数据库初始化失败")
	}

	if len(EasyGorm.Config.Models) > 0 {
		err = Migrate(EasyGorm.Config.Models)
		if err != nil {
			return errors.New(fmt.Sprintf("AutoMigrate err: : %+v", err))
		}
	}

	// 没有 CasbinPath 不使用 casbin
	if EasyGorm.Casbin == nil {
		return errors.New("casbin 设置不能为空")
	}

	err = EasyGorm.getEnforcer()
	if err != nil {
		return errors.New(fmt.Sprintf("getEnforcer err: : %+v", err))
	}

	if EasyGorm.Enforcer == nil {
		return errors.New("casbin 初始化失败")
	}

	return nil
}

// Migrate 迁移数据表
func Migrate(models []interface{}) error {
	err := EasyGorm.DB.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}

// getGormDb
func (db *DBServer) getGormDb() error {
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
	db.DB, err = gorm.Open(dialector, db.Config.GormConfig)
	if err != nil {
		return err
	}

	err = db.DB.Use(
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

	db.DB.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: false})
	return nil
}

// getEnforcer
func (db *DBServer) getEnforcer() error {
	c, err := gormadapter.NewAdapterByDBUseTableName(db.DB, db.Config.Casbin.Prefix, "casbin_rule") // Your driver and data source.
	if err != nil {
		return err
	}

	db.Enforcer, err = casbin.NewEnforcer(db.Config.Casbin.Path, c)
	if err != nil {
		return err
	}

	err = db.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}

	return nil
}
