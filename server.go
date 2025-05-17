package admin

import (
	"fmt"
	"log"
	"time"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/mattn/go-colorable"
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
	conf     *conf.Conf
	db       *gorm.DB
	enforcer *casbin.Enforcer
	engine   *gin.Engine
	iroutes  *gin.IRoutes

	validate *Validator

	m     *gormigrate.Gormigrate
	items []*gormigrate.Migration

	permRoutes  []*Router
	otherRoutes []*Router
}

// gormDb
func gormDb(m *conf.Mysql) (*gorm.DB, error) {
	if m == nil {
		return nil, e.ErrConfigInvalid
	}
	if m.DbName == "" {
		return nil, e.ErrDbTableNameEmpty
	}
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
	app := gin.Default()
	if c.System.Tls {
		app.Use(LoadTls())
	}
	app.Use(c.CorsConf.Cors())
	// registerValidation()
	gin.DefaultWriter = colorable.NewColorableStdout()
	c.SetDefaultAddrAndTimeFormat()
	db, err := gormDb(c.Mysql)
	if err != nil {
		return nil, err
	}

	auth, err := c.GetEnforcer(db)
	if err != nil {
		return nil, err
	}

	ws := &WebServe{
		conf:        c,
		engine:      app,
		enforcer:    auth,
		db:          db,
		permRoutes:  []*Router{},
		otherRoutes: []*Router{},
	}
	if err := ws.Migrate(); err != nil {
		return nil, err
	}

	switch c.Locale {
	case "en":
		ws.validate = newEn()
	case "zh":
		ws.validate = newZh()
	default:
		ws.validate = newZh()
	}

	ws.engine.Use(limit.MaxAllowed(50))

	return ws, nil
}

// Engine return *gin.Engine
func (ws *WebServe) Engine() *gin.Engine {
	return ws.engine
}

func (ws *WebServe) IRoutes() *gin.IRoutes {
	return ws.iroutes
}

// Config
func (ws *WebServe) Config() *conf.Conf {
	return ws.conf
}

// SystemAddr
func (ws *WebServe) SystemAddr() string {
	return ws.conf.System.Addr
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
	if ws.engine == nil {
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

	ws.groupRouters()

	systemAddr := ws.SystemAddr()
	s := run(systemAddr, ws.engine)
	time.Sleep(10 * time.Microsecond)

	log.Printf("listen on: http://%s\n", systemAddr)

	s.ListenAndServe()
}
