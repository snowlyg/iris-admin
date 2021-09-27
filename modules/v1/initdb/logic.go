package initdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/modules/v1/perm"
	"github.com/snowlyg/iris-admin/modules/v1/role"
	"github.com/snowlyg/iris-admin/modules/v1/user"
	"github.com/snowlyg/iris-admin/server/config"
	"github.com/snowlyg/iris-admin/server/database"
	myviper "github.com/snowlyg/iris-admin/server/viper"
	myzap "github.com/snowlyg/iris-admin/server/zap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ErrViperEmpty = errors.New("配置服务未初始化")
)

// writeConfig 写入配置文件
func writeConfig(viper *viper.Viper, conf config.Config) error {
	cs := str.StructToMap(config.CONFIG)
	for k, v := range cs {
		viper.Set(k, v)
	}
	return viper.WriteConfig()
}

// 回滚配置
func refreshConfig(viper *viper.Viper, conf config.Config) error {
	err := writeConfig(viper, conf)
	if err != nil {
		myzap.ZAPLOG.Error("还原配置文件设置错误", zap.String("refreshConfig(g.VIPER)", err.Error()))
		return err
	}
	return nil
}

// createTable 创建数据库(mysql)
func createTable(dsn string, driver string, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}

// initDB 初始化数据
func initDB(InitDBFunctions ...database.InitDBFunc) error {
	for _, v := range InitDBFunctions {
		err := v.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

// InitDB 创建数据库并初始化
func InitDB(req *Request) error {
	defaultConfig := config.CONFIG
	if myviper.VIPER == nil {
		myzap.ZAPLOG.Error("初始化错误", zap.String("InitDB", ErrViperEmpty.Error()))
		return ErrViperEmpty
	}

	level := req.Level
	if level == "" {
		level = "release"
	}
	addr := req.Addr
	if addr == "" {
		addr = "127.0.0.1:8085"
	}

	config.CONFIG.System.CacheType = req.CacheType
	config.CONFIG.System.Level = level
	config.CONFIG.System.Addr = addr
	config.CONFIG.System.DbType = req.SqlType

	if config.CONFIG.System.CacheType == "redis" {
		config.CONFIG.Redis = config.Redis{
			DB:       req.Cache.DB,
			Addr:     fmt.Sprintf("%s:%s", req.Cache.Host, req.Cache.Port),
			Password: req.Cache.Password,
		}
	}

	if req.Sql.Host == "" {
		req.Sql.Host = "127.0.0.1"
	}

	if req.Sql.Port == "" {
		req.Sql.Port = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", req.Sql.UserName, req.Sql.Password, req.Sql.Host, req.Sql.Port)
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", req.Sql.DBName)

	if err := createTable(dsn, "mysql", createSql); err != nil {
		return err
	}

	myzap.ZAPLOG.Info("新建数据库", zap.String("库名", req.Sql.DBName))

	config.CONFIG.Mysql.Path = fmt.Sprintf("%s:%s", req.Sql.Host, req.Sql.Port)
	config.CONFIG.Mysql.Dbname = req.Sql.DBName
	config.CONFIG.Mysql.Username = req.Sql.UserName
	config.CONFIG.Mysql.Password = req.Sql.Password
	config.CONFIG.Mysql.LogMode = req.Sql.LogMode

	m := config.CONFIG.Mysql
	if m.Dbname == "" {
		myzap.ZAPLOG.Error("缺少数据库参数")
		return errors.New("缺少数据库参数")
	}

	if err := writeConfig(myviper.VIPER, config.CONFIG); err != nil {
		myzap.ZAPLOG.Error("更新配置文件错误", zap.String("writeConfig(g.VIPER)", err.Error()))
	}

	if database.Instance() == nil {
		myzap.ZAPLOG.Error("数据库初始化错误")
		refreshConfig(myviper.VIPER, defaultConfig)
		return errors.New("数据库初始化错误")
	}

	err := database.Instance().AutoMigrate(
		&middleware.Oplog{},
		&perm.Permission{},
		&role.Role{},
		&user.User{},
	)
	if err != nil {
		myzap.ZAPLOG.Error("迁移数据表错误", zap.String("错误:", err.Error()))
		refreshConfig(myviper.VIPER, defaultConfig)
		return err
	}

	err = initDB(
		perm.Source,
		role.Source,
		user.Source,
	)
	if err != nil {
		myzap.ZAPLOG.Error("填充数据错误", zap.String("错误:", err.Error()))
		refreshConfig(myviper.VIPER, defaultConfig)
		return err
	}
	return nil
}
