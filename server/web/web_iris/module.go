package web_iris

// // InitParty 初始化模块
// // 分离操作和数据,使用命令行处理项目初始化
// func InitParty() module.WebModule {
// 	handler := func(index iris.Party) {
// 		index.Post("/initdb", InitDB)
// 		index.Get("/checkdb", Check)
// 	}
// 	return module.NewModule("/init", handler)
// }

// // InitDB 初始化项目接口
// func InitDB(ctx iris.Context) {
// 	req := &Request{}
// 	if err := req.Request(ctx); err != nil {
// 		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
// 		return
// 	}
// 	err := InitWeb(req)
// 	if err != nil {
// 		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: orm.SystemErr.Msg})
// 		return
// 	}
// 	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: nil, Msg: orm.NoErr.Msg})
// }

// // Check 检测是否需要初始化项目
// func Check(ctx iris.Context) {
// 	if database.Instance() == nil {
// 		ctx.JSON(orm.Response{Code: orm.NeedInitErr.Code, Data: iris.Map{
// 			"needInit": true,
// 		}, Msg: str.Join(orm.NeedInitErr.Msg, ":数据库初始化失败")})
// 		return
// 	} else if config.CONFIG.System.CacheType == "redis" && cache.Instance() == nil {
// 		ctx.JSON(orm.Response{Code: orm.NeedInitErr.Code, Data: iris.Map{
// 			"needInit": true,
// 		}, Msg: str.Join(orm.NeedInitErr.Msg, ":缓存驱动初始化失败")})
// 		return
// 	}
// 	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: iris.Map{
// 		"needInit": false,
// 	}, Msg: orm.NoErr.Msg})
// }

// type Request struct {
// 	Sql       Sql    `json:"sql"`
// 	SqlType   string `json:"sqlType" validate:"required"`
// 	Cache     Cache  `json:"cache"`
// 	CacheType string `json:"cacheType"  validate:"required"`
// 	Level     string `json:"level"` // debug,release,test
// 	Addr      string `json:"addr"`
// }

// func (req *Request) Request(ctx iris.Context) error {
// 	if err := ctx.ReadJSON(req); err != nil {
// 		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
// 		return orm.ErrParamValidate
// 	}
// 	return nil
// }

// type Sql struct {
// 	Host     string `json:"host"  validate:"required"`
// 	Port     string `json:"port"  validate:"required"`
// 	UserName string `json:"userName" validate:"required"`
// 	Password string `json:"password"  validate:"required"`
// 	DBName   string `json:"dbName" validate:"required"`
// 	LogMode  bool   `json:"logMode"`
// }

// type Cache struct {
// 	Host     string `json:"host"  validate:"required"`
// 	Port     string `json:"port"  validate:"required"`
// 	Password string `json:"password"`
// 	PoolSize int    `json:"poolSize"`
// 	DB       int    `json:"db"`
// }

// var (
// 	ErrViperEmpty = errors.New("配置服务未初始化")
// )

// // writeConfig 写入配置文件
// func writeConfig(viper *viper.Viper, conf config.Config) error {
// 	cs := str.StructToMap(config.CONFIG)
// 	for k, v := range cs {
// 		viper.Set(k, v)
// 	}
// 	return viper.WriteConfig()
// }

// // 回滚配置
// func refreshConfig(viper *viper.Viper, conf config.Config) error {
// 	err := writeConfig(viper, conf)
// 	if err != nil {
// 		zap_server.ZAPLOG.Error("还原配置文件设置错误", zap.String("refreshConfig(g.VIPER)", err.Error()))
// 		return err
// 	}
// 	return nil
// }

// // createTable 创建数据库(mysql)
// func createTable(dsn string, driver string, createSql string) error {
// 	db, err := sql.Open(driver, dsn)
// 	if err != nil {
// 		return err
// 	}
// 	defer func(db *sql.DB) {
// 		_ = db.Close()
// 	}(db)
// 	if err = db.Ping(); err != nil {
// 		return err
// 	}
// 	_, err = db.Exec(createSql)
// 	return err
// }

// // initDB 初始化数据
// func initDB(InitDBFunctions ...module.InitSourceFunc) error {
// 	for _, v := range InitDBFunctions {
// 		err := v.Init()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // InitWeb 创建数据库并初始化
// func InitWeb(req *Request) error {
// 	defaultConfig := config.CONFIG
// 	if viper_server.VIPER == nil {
// 		zap_server.ZAPLOG.Error("初始化错误", zap.String("InitWeb()", ErrViperEmpty.Error()))
// 		return ErrViperEmpty
// 	}

// 	level := req.Level
// 	if level == "" {
// 		level = "release"
// 	}
// 	addr := req.Addr
// 	if addr == "" {
// 		addr = "127.0.0.1:8085"
// 	}

// 	config.CONFIG.System.CacheType = req.CacheType
// 	config.CONFIG.System.Level = level
// 	config.CONFIG.System.Addr = addr
// 	config.CONFIG.System.DbType = req.SqlType

// 	if config.CONFIG.System.CacheType == "redis" {
// 		config.CONFIG.Redis = config.Redis{
// 			DB:       req.Cache.DB,
// 			Addr:     fmt.Sprintf("%s:%s", req.Cache.Host, req.Cache.Port),
// 			Password: req.Cache.Password,
// 		}
// 	}

// 	if req.Sql.Host == "" {
// 		req.Sql.Host = "127.0.0.1"
// 	}

// 	if req.Sql.Port == "" {
// 		req.Sql.Port = "3306"
// 	}

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", req.Sql.UserName, req.Sql.Password, req.Sql.Host, req.Sql.Port)
// 	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", req.Sql.DBName)

// 	if err := createTable(dsn, "mysql", createSql); err != nil {
// 		return err
// 	}

// 	zap_server.ZAPLOG.Info("新建数据库", zap.String("库名", req.Sql.DBName))

// 	config.CONFIG.Mysql.Path = fmt.Sprintf("%s:%s", req.Sql.Host, req.Sql.Port)
// 	config.CONFIG.Mysql.Dbname = req.Sql.DBName
// 	config.CONFIG.Mysql.Username = req.Sql.UserName
// 	config.CONFIG.Mysql.Password = req.Sql.Password
// 	config.CONFIG.Mysql.LogMode = req.Sql.LogMode

// 	m := config.CONFIG.Mysql
// 	if m.Dbname == "" {
// 		zap_server.ZAPLOG.Error("缺少数据库参数")
// 		return errors.New("缺少数据库参数")
// 	}

// 	if err := writeConfig(viper_server.VIPER, config.CONFIG); err != nil {
// 		zap_server.ZAPLOG.Error("更新配置文件错误", zap.String("writeConfig(g.VIPER)", err.Error()))
// 	}

// 	if database.Instance() == nil {
// 		zap_server.ZAPLOG.Error("数据库初始化错误")
// 		refreshConfig(viper_server.VIPER, defaultConfig)
// 		return errors.New("数据库初始化错误")
// 	}

// 	// TODO: 通过全局变量获取模型
// 	err := database.Instance().AutoMigrate(
// 		&middleware.Oplog{},
// 		&perm.Permission{},
// 		&role.Role{},
// 		&user.User{},
// 	)
// 	if err != nil {
// 		zap_server.ZAPLOG.Error("迁移数据表错误", zap.String("错误:", err.Error()))
// 		refreshConfig(viper_server.VIPER, defaultConfig)
// 		return err
// 	}

// 	// TODO: 通过全局变量填充数据
// 	err = initDB(
// 		perm.Source,
// 		role.Source,
// 		user.Source,
// 	)
// 	if err != nil {
// 		zap_server.ZAPLOG.Error("填充数据错误", zap.String("错误:", err.Error()))
// 		refreshConfig(viper_server.VIPER, defaultConfig)
// 		return err
// 	}
// 	return nil
// }
