package admin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

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
	serveMu   sync.RWMutex
	closeOnce sync.Once
	closeErr  error
	conf      *conf.Conf
	db        *gorm.DB
	enforcer  *casbin.Enforcer
	engine    *gin.Engine
	iroutes   *gin.IRoutes

	validate *Validator

	m     *gormigrate.Gormigrate
	items []*gormigrate.Migration

	permRoutes  []*Router
	otherRoutes []*Router
}

var (
	ErrServerNotInitialized    = errors.New("server_not_initialized")
	ErrShutdownContextRequired = errors.New("shutdown_context_required")
)

type serverDependencies struct {
	ensureRbac func(*conf.Conf) error
	openDB     func(*conf.Mysql) (*gorm.DB, error)
	getAuth    func(*conf.Conf, *gorm.DB) (*casbin.Enforcer, error)
	migrate    func(*WebServe) error
	closeDB    func(*gorm.DB) error
}

func defaultServerDependencies() serverDependencies {
	return serverDependencies{
		ensureRbac: (*conf.Conf).EnsureRbacModel,
		openDB:     gormDb,
		getAuth:    (*conf.Conf).GetEnforcer,
		migrate:    (*WebServe).Migrate,
		closeDB:    closeGormDB,
	}
}

// gormDb
func gormDb(m *conf.Mysql) (*gorm.DB, error) {
	return gormDbWithOpen(m, gorm.Open)
}

func gormDbWithOpen(
	m *conf.Mysql,
	open func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error),
) (*gorm.DB, error) {
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
	if db, err := open(mysql.New(mysqlConfig)); err != nil {
		log.Printf("open mysql %s failed\n", mysqlLogTarget(m))
		return nil, err
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		if err := sqlDB.Ping(); err != nil {
			log.Printf("ping mysql %s failed\n", mysqlLogTarget(m))
			return nil, err
		}
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db, nil
	}
}

func mysqlLogTarget(m *conf.Mysql) string {
	return fmt.Sprintf("address=%q database=%q", m.Path, m.DbName)
}

func closeGormDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// NewServe
func NewServe(c *conf.Conf) (*WebServe, error) {
	return newServe(c, defaultServerDependencies())
}

func newServe(c *conf.Conf, dependencies serverDependencies) (*WebServe, error) {
	if c != nil {
		c.SetDefaultAddrAndTimeFormat()
	}
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}
	if err := dependencies.ensureRbac(c); err != nil {
		return nil, err
	}
	gin.SetMode(c.System.GinMode)
	app := gin.Default()
	if c.System.Tls {
		app.Use(LoadTls())
	}
	app.Use(c.CorsConf.Cors())
	// registerValidation()
	gin.DefaultWriter = colorable.NewColorableStdout()
	db, err := dependencies.openDB(c.Mysql)
	if err != nil {
		return nil, err
	}

	auth, err := dependencies.getAuth(c, db)
	if err != nil {
		return nil, errors.Join(err, dependencies.closeDB(db))
	}

	ws := &WebServe{
		conf:        c,
		engine:      app,
		enforcer:    auth,
		db:          db,
		permRoutes:  []*Router{},
		otherRoutes: []*Router{},
	}
	if err := dependencies.migrate(ws); err != nil {
		return nil, errors.Join(err, dependencies.closeDB(db))
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

// Run starts the service and logs startup or serving errors.
func (ws *WebServe) Run() {
	if err := ws.RunE(); err != nil {
		log.Printf("iris-admin: server stopped with error: %v\n", err)
	}
}

// RunE starts the service and returns startup or serving errors.
func (ws *WebServe) RunE() error {
	if ws == nil || ws.engine == nil || ws.conf == nil {
		return ErrServerNotInitialized
	}
	ws.groupRouters()

	ws.serveMu.Lock()
	if ws.serve == nil {
		ws.serve = run(ws.SystemAddr(), ws.engine)
	}
	server := ws.serve
	ws.serveMu.Unlock()

	log.Printf("listen on: http://%s\n", ws.SystemAddr())
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown gracefully stops the HTTP server and closes the database.
func (ws *WebServe) Shutdown(ctx context.Context) error {
	if ctx == nil {
		return ErrShutdownContextRequired
	}
	if ws == nil {
		return ErrServerNotInitialized
	}

	ws.serveMu.RLock()
	server := ws.serve
	ws.serveMu.RUnlock()

	var shutdownErr error
	if server != nil {
		shutdownErr = server.Shutdown(ctx)
	}
	ws.closeOnce.Do(func() {
		ws.closeErr = closeGormDB(ws.db)
	})
	return errors.Join(shutdownErr, ws.closeErr)
}
