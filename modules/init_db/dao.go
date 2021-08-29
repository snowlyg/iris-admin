package init_db

import (
	"database/sql"
	"fmt"
	"os/user"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/config"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/modules/perm"
	"github.com/snowlyg/iris-admin/modules/role"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/module"
	"github.com/snowlyg/multi"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	baseSystem = config.System{
		CacheType: "",
		Level:     "debug",
		Addr:      ":80",
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
	env := req.Level
	if env == "" {
		env = "pro"
	}
	addr := req.Addr
	if addr == "" {
		addr = ":80"
	}

	g.CONFIG.System = config.System{
		CacheType: req.CacheType,
		Level:     level,
		Addr:      addr,
		DbType:    req.SqlType,
	}

	if g.CONFIG.System.CacheType == "redis" {
		g.CONFIG.Redis = config.Redis{
			DB:       0,
			Addr:     fmt.Sprintf("%s:%s", req.Cache.Host, req.Cache.Port),
			Password: req.Cache.Password,
		}
		cache.Init() // redis缓存
		err := multi.InitDriver(
			&multi.Config{
				DriverType:      g.CONFIG.System.CacheType,
				UniversalClient: g.CACHE},
		)
		if err != nil {
			g.ZAPLOG.Error("初始化缓存驱动:", zap.Any("err", err))
			if err := refreshConfig(g.VIPER); err != nil {
				g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误", err.Error()))
			}
			return fmt.Errorf("初始化缓存驱动失败 %w", err)
		}
		if multi.AuthDriver == nil {
			if err := refreshConfig(g.VIPER); err != nil {
				g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误", err.Error()))
			}
			return nil
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
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误", err.Error()))
		}
		return err
	}

	g.CONFIG.Mysql = config.Mysql{
		Path:     fmt.Sprintf("%s:%s", req.Sql.Host, req.Sql.Port),
		Dbname:   req.Sql.DBName,
		Username: req.Sql.UserName,
		Password: req.Sql.Password,
		Config:   "charset=utf8mb4&parseTime=True&loc=Local",
	}

	m := g.CONFIG.Mysql
	if m.Dbname == "" {
		if err := refreshConfig(g.VIPER); err != nil {
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误", err.Error()))
		}
		return nil
	}

	if database.Instance() == nil {
		if err := refreshConfig(g.VIPER); err != nil {
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误", err.Error()))
		}
		return nil
	}

	err := database.Instance().AutoMigrate(
		&middleware.Oplog{},
		&perm.Permission{},
		&role.Role{},
		&user.User{},
	)
	if err != nil {
		if err := refreshConfig(g.VIPER); err != nil {
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误", err.Error()))
		}
		return err
	}

	err = initDB()
	if err != nil {
		if err := refreshConfig(g.VIPER); err != nil {
			g.ZAPLOG.Error("还原配置文件设置错误", zap.String("错误", err.Error()))
		}
		return err
	}
	if err := writeConfig(g.VIPER); err != nil {
		g.ZAPLOG.Error("更新配置文件错误", zap.String("错误", err.Error()))
	}
	return nil
}
