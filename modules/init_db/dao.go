package init_db

import (
	"fmt"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/config"
	"github.com/snowlyg/iris-admin/g"
)

func writeConfig() {
	g.VIPER.Set("verbose", true)
	g.VIPER.WriteConfig()
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
	BaseSystem := config.System{
		CacheType: req.CacheType,
		Level:     level,
		Addr:      addr,
		DbType:    req.SqlType,
	}

	for _, conf := range str.StructToMap(BaseSystem) {
		fmt.Printf("--- %+v\n", conf)
	}
	// if err := WriteCacheTypeConfig(g.VIPER, BaseSystem); err != nil {
	// 	return err
	// }
	// if BaseSystem.CacheType == "redis" {
	// 	BaseCache := config.Redis{
	// 		DB:       0,
	// 		Addr:     fmt.Sprintf("%s:%s", conf.Cache.Host, conf.Cache.Port),
	// 		Password: conf.Cache.Password,
	// 	}
	// 	if err := WriteRedisConfig(g.TENANCY_VP, BaseCache); err != nil {
	// 		return err
	// 	}
	// 	g.TENANCY_CACHE = cache.Cache() // redis缓存
	// 	err := multi.InitDriver(&multi.Config{
	// 		DriverType:      g.TENANCY_CONFIG.System.CacheType,
	// 		UniversalClient: g.TENANCY_CACHE})
	// 	if err != nil {
	// 		g.TENANCY_LOG.Error("初始化缓存驱动:", zap.Any("err", err))
	// 		return fmt.Errorf("初始化缓存驱动失败 %w", err)
	// 	}
	// 	if multi.AuthDriver == nil {

	// 	}
	// }

	// BaseMysql := config.Mysql{
	// 	Path:     "",
	// 	Dbname:   "",
	// 	Username: "",
	// 	Password: "",
	// 	Config:   "charset=utf8mb4&parseTime=True&loc=Local",
	// }

	// if conf.Sql.Host == "" {
	// 	conf.Sql.Host = "127.0.0.1"
	// }

	// if conf.Sql.Port == "" {
	// 	conf.Sql.Port = "3306"
	// }
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", conf.Sql.UserName, conf.Sql.Password, conf.Sql.Host, conf.Sql.Port)
	// createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", conf.Sql.DBName)

	// if err := createTable(dsn, "mysql", createSql); err != nil {
	// 	return err
	// }

	// MysqlConfig := config.Mysql{
	// 	Path:     fmt.Sprintf("%s:%s", conf.Sql.Host, conf.Sql.Port),
	// 	Dbname:   conf.Sql.DBName,
	// 	Username: conf.Sql.UserName,
	// 	Password: conf.Sql.Password,
	// 	Config:   "charset=utf8mb4&parseTime=True&loc=Local",
	// }

	// if err := WriteConfig(g.TENANCY_VP, MysqlConfig); err != nil {
	// 	return err
	// }
	// m := g.TENANCY_CONFIG.Mysql
	// if m.Dbname == "" {
	// 	return nil
	// }

	// linkDns := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	// mysqlConfig := mysql.Config{
	// 	DSN:                       linkDns, // DSN data source name
	// 	DefaultStringSize:         191,     // string 类型字段的默认长度
	// 	DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	// 	DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	// 	DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	// 	SkipInitializeWithVersion: false,   // 根据版本自动配置
	// }
	// if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
	// 	_ = WriteConfig(g.TENANCY_VP, BaseMysql)
	// 	return nil
	// } else {
	// 	sqlDB, _ := db.DB()
	// 	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	// 	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	// 	g.TENANCY_DB = db
	// }

	// err := g.TENANCY_DB.AutoMigrate(
	// 	model.SysUser{},
	// 	model.AdminInfo{},
	// 	model.GeneralInfo{},
	// 	model.SysAuthority{},
	// 	model.SysApi{},
	// 	model.SysBaseMenu{},
	// 	model.SysRegion{},
	// 	model.SysOperationRecord{},
	// 	model.SysTenancy{},
	// 	model.SysMini{},
	// 	model.SysConfig{},
	// 	model.SysConfigCategory{},
	// 	model.SysConfigValue{},
	// 	model.SysBrandCategory{},
	// 	model.SysBrand{},
	// 	model.Patient{},

	// 	model.TenancyMedia{},
	// 	model.ProductCategory{},
	// 	model.AttrTemplate{},
	// 	model.Product{},
	// 	model.ProductProductCate{},
	// 	model.ProductContent{},
	// 	model.ProductAttrValue{},
	// 	model.ProductReply{},
	// 	model.ShippingTemplate{},
	// 	model.ShippingTemplateFree{},
	// 	model.ShippingTemplateRegion{},
	// 	model.ShippingTemplateUndelivery{},

	// 	model.Cart{},
	// 	model.Express{},
	// 	model.Order{},
	// 	model.OrderStatus{},
	// 	model.OrderReceipt{},
	// 	model.OrderProduct{},
	// 	model.GroupOrder{},
	// 	model.CartOrder{},

	// 	model.RefundOrder{},
	// 	model.RefundProduct{},
	// 	model.RefundStatus{},

	// 	model.UserAddress{},
	// 	model.UserReceipt{},
	// 	model.UserBill{},
	// 	model.UserExtract{},
	// 	model.UserGroup{},
	// 	model.UserLabel{},
	// 	model.UserUserLabel{},
	// 	model.LabelRule{},
	// 	model.UserMerchant{},
	// 	model.UserRecharge{},
	// 	model.UserRelation{},
	// 	model.UserVisit{},
	// 	model.Mqtt{},
	// 	model.MqttRecord{},
	// )
	// if err != nil {
	// 	_ = WriteConfig(g.TENANCY_VP, BaseMysql)
	// 	return err
	// }
	// err = initDB(
	// 	source.Admin,
	// 	source.Api,
	// 	source.AuthorityMenu,
	// 	source.Authority,
	// 	source.AuthoritiesMenus,
	// 	source.Casbin,
	// 	source.DataAuthorities,
	// 	source.BaseMenu,
	// 	source.Region,
	// 	source.Config,
	// 	source.SysConfigCategory,
	// 	source.SysConfigValue,
	// )
	// if err != nil {
	// 	_ = WriteConfig(g.TENANCY_VP, BaseMysql)
	// 	return err
	// }
	return nil
}
