package initdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/config"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/modules/v1/perm"
	"github.com/snowlyg/iris-admin/modules/v1/role"
	"github.com/snowlyg/iris-admin/modules/v1/user"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/module"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	baseSystem = config.System{
		CacheType: "",
		Level:     "debug",
		Addr:      "127.0.0.1:80",
		DbType:    "",
	}
	baseCache = config.Redis{
		DB:       0,
		Addr:     "",
		Password: "",
	}
	baseMysql = config.Mysql{
		Path:     "",
		Dbname:   "",
		Username: "",
		Password: "",
		Config:   "charset=utf8mb4&parseTime=True&loc=Local",
		LogMode:  true,
	}
)

func writeConfig(viper *viper.Viper) error {
	cs := str.StructToMap(g.CONFIG)
	for k, v := range cs {
		viper.Set(k, v)
	}
	return viper.WriteConfig()
}

func refreshConfig(viper *viper.Viper) error {
	g.CONFIG.System = baseSystem
	g.CONFIG.Redis = baseCache
	g.CONFIG.Mysql = baseMysql
	cs := str.StructToMap(g.CONFIG)
	for k, v := range cs {
		viper.Set(k, v)
	}
	return viper.WriteConfig()
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

func initDB(InitDBFunctions ...module.InitDBFunc) error {
	for _, v := range InitDBFunctions {
		err := v.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

// InitDB 创建数据库并初始化
func InitDB(req Request) error {
	level := req.Level
	if level == "" {
		level = "release"
	}
	addr := req.Addr
	if addr == "" {
		addr = "127.0.0.1:8085"
	}

	g.CONFIG.System = config.System{
		CacheType: req.CacheType,
		Level:     level,
		Addr:      addr,
		DbType:    req.SqlType,
	}

	if g.CONFIG.System.CacheType == "redis" {
		g.CONFIG.Redis = config.Redis{
			DB:       req.Cache.DB,
			Addr:     fmt.Sprintf("%s:%s", req.Cache.Host, req.Cache.Port),
			Password: req.Cache.Password,
		}
		err := cache.Init() // redis缓存
		if err != nil {
			g.ZAPLOG.Error("认证驱动初始化错误", zap.String("错误:", err.Error()))
			if err := refreshConfig(g.VIPER); err != nil {
				g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误:", err.Error()))
			}
			return err
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
		if err := refreshConfig(g.VIPER); err != nil {
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误:", err.Error()))
		}
		return err
	}

	g.ZAPLOG.Info("新建数据库", zap.String("库名", req.Sql.DBName))

	g.CONFIG.Mysql = config.Mysql{
		Path:     fmt.Sprintf("%s:%s", req.Sql.Host, req.Sql.Port),
		Dbname:   req.Sql.DBName,
		Username: req.Sql.UserName,
		Password: req.Sql.Password,
		Config:   "charset=utf8mb4&parseTime=True&loc=Local",
		LogMode:  req.Sql.LogMode,
	}

	m := g.CONFIG.Mysql
	if m.Dbname == "" {
		if err := refreshConfig(g.VIPER); err != nil {
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误:", err.Error()))
		}
		g.ZAPLOG.Error("缺少数据库参数")
		return errors.New("缺少数据库参数")
	}

	if err := writeConfig(g.VIPER); err != nil {
		g.ZAPLOG.Error("更新配置文件错误", zap.String("错误:", err.Error()))
	}

	if database.Instance() == nil {
		if err := refreshConfig(g.VIPER); err != nil {
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误:", err.Error()))
		}
		g.ZAPLOG.Error("数据库初始化错误")
		return errors.New("数据库初始化错误")
	}

	err := database.Instance().AutoMigrate(
		&middleware.Oplog{},
		&perm.Permission{},
		&role.Role{},
		&user.User{},
	)
	if err != nil {
		g.ZAPLOG.Error("迁移数据表错误", zap.String("错误:", err.Error()))
		if err := refreshConfig(g.VIPER); err != nil {
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误:", err.Error()))
		}
		return err
	}

	err = initDB(
		perm.Source,
		role.Source,
		user.Source,
	)
	if err != nil {
		g.ZAPLOG.Error("填充数据错误", zap.String("错误:", err.Error()))
		if err := refreshConfig(g.VIPER); err != nil {
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误:", err.Error()))
		}
		return err
	}
	return nil
}
