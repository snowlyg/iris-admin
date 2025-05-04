package admin

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/conf"
	"github.com/snowlyg/iris-admin/e"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Status
const (
	StatusUnknown int = iota
	StatusTrue
	StatusFalse
)

type WebServe struct {
	serve
	db       *gorm.DB
	enforcer *casbin.Enforcer
	engine   *gin.Engine
	conf     *conf.Conf
	validate *Validator
	m        *Migrate
}

// getEnforcer get casbin.Enforcer
func getEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	c, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule") // Your driver and data source.
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(filepath.Join(dir.GetCurrentAbPath(), conf.CasbinName), c)
	if err != nil {
		return nil, err
	}

	if err = enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return enforcer, nil
}

// gormDb
func gormDb(m *conf.Mysql) (*gorm.DB, error) {
	if m == nil {
		return nil, e.ErrConfigInvalid
	}
	if m.DbName == "" {
		return nil, e.ErrDbTableNameEmpty
	}
	// if err := createTable(conf.BaseDsn(), "mysql", conf.DbName); err != nil {
	// 	return nil, fmt.Errorf("create database %s is fail:%w", conf.DbName, err)
	// }
	mysqlConfig := mysql.Config{
		DSN:               m.Dsn(),
		DefaultStringSize: 191,
		// DisableDatetimePrecision:  true,
		// DontSupportRenameIndex:    true,
		// DontSupportRenameColumn:   true,
		// SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig)); err != nil {
		fmt.Printf("open mysql[%s] is fail:%v\n", m.Dsn(), err)
		return nil, err
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		if err := sqlDB.Ping(); err != nil {
			log.Printf("ping mysql[%s] is fail:%v\n", m.Dsn(), err)
			return nil, err
		}
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db, nil
	}
}

// NewServe
func NewServe(c *conf.Conf) (*WebServe, error) {
	gin.SetMode(c.System.GinMode)
	log.Printf("set mode:%s\n", c.System.GinMode)
	app := gin.Default()
	if c.System.Tls {
		app.Use(LoadTls())
		log.Printf("use tls\n")
	}
	app.Use(c.CorsConf.Cors())
	log.Printf("use cors\n")
	// registerValidation()
	gin.DefaultWriter = colorable.NewColorableStdout()
	c.SetDefaultAddrAndTimeFormat()
	log.Printf("set default addr:%s and time format:%s\n", c.System.Addr, c.System.TimeFormat)
	db, err := gormDb(c.Mysql)
	if err != nil {
		return nil, err
	}
	log.Printf("init gorm database, dsn: %s\n", c.Mysql.Dsn())
	auth, err := getEnforcer(db)
	if err != nil {
		return nil, err
	}
	log.Printf("init casbin :%s\n", conf.CasbinName)
	webServe := &WebServe{
		conf:     c,
		engine:   app,
		enforcer: auth,
		db:       db,
		m: &Migrate{
			db:    db,
			items: nil,
			seeds: nil,
		},
	}
	switch c.Locale {
	case "en":
		webServe.validate = newEn()
	case "zh":
		webServe.validate = newZh()
	default:
		webServe.validate = newZh()
	}
	return webServe, nil
}

// Engine return *gin.Engine
func (ws *WebServe) Engine() *gin.Engine {
	return ws.engine
}

// Config
func (ws *WebServe) Config() *conf.Conf {
	return ws.conf
}

// Auth
func (ws *WebServe) Auth() *casbin.Enforcer {
	return ws.enforcer
}

// DB
func (ws *WebServe) DB() *gorm.DB {
	return ws.db
}

// // Deprecated: use nginx or apache instead.
// func (ws *WebServe) AddWebStatic(staticAbsPath, webPrefix string, paths ...string) {
// }

// // Deprecated: use nginx or apache instead.
// func (ws *WebServe) AddUploadStatic(webPrefix, staticAbsPath string) {
// }

// Run
func (ws *WebServe) Run() {
	if ws.Engine() == nil {
		panic("init engine please")
	}

	// ws.Engine().NoRoute(func(ctx *gin.Context) {
	// 	// excepte for /v0 /v1 and so on
	// 	reg := `^/v[0-9]+$|^(/v[0-9]+)/`
	// 	ok, _ := regexp.MatchString(reg, ctx.Request.RequestURI)
	// 	if ok {
	// 		ctx.Writer.WriteHeader(http.StatusNotFound)
	// 		ctx.Writer.Flush()
	// 		return
	// 	}

	// 	var indexFile []byte
	// 	for _, wp := range ws.statics {
	// 		// match /admin or /admin/***
	// 		reg := str.Join("^", wp.Prefix, "$|^(", wp.Prefix, ")/")
	// 		ok, err := regexp.MatchString(reg, ctx.Request.RequestURI)
	// 		if err != nil || !ok {
	// 			continue
	// 		}
	// 		indexFile = wp.IndexFile
	// 	}

	// 	ctx.Writer.WriteHeader(http.StatusOK)
	// 	ctx.Writer.Write(indexFile)

	// 	ctx.Writer.Header().Add("Accept", "text/html")
	// 	ctx.Writer.Flush()
	// })
	s := run(ws.Config().System.Addr, ws.engine)
	time.Sleep(10 * time.Microsecond)
	log.Printf("listen on: http://%s\n", ws.Config().System.Addr)

	s.ListenAndServe()
}
