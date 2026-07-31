package admin

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/conf"
	"github.com/snowlyg/iris-admin/e"
	"gorm.io/gorm"
)

type fakeServe struct {
	listenErr     error
	shutdownErr   error
	listenCalls   int
	shutdownCalls int
}

func (server *fakeServe) ListenAndServe() error {
	server.listenCalls++
	return server.listenErr
}

func (server *fakeServe) Shutdown(context.Context) error {
	server.shutdownCalls++
	return server.shutdownErr
}

func TestNewServeValidatesConfigBeforeStartup(t *testing.T) {
	if _, err := NewServe(nil); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("NewServe(nil) error = %v, want %v", err, e.ErrConfigInvalid)
	}
	if _, err := NewServe(&conf.Conf{}); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("NewServe(empty) error = %v, want %v", err, e.ErrConfigInvalid)
	}
}

func TestNewServeHandlesStartupDependencies(t *testing.T) {
	config := conf.NewConf()
	database := &gorm.DB{}
	enforcer := &casbin.Enforcer{}

	t.Run("rbac model", func(t *testing.T) {
		rbacErr := errors.New("rbac model failed")
		dependencies := successfulServerDependencies(database, enforcer)
		dependencies.ensureRbac = func(*conf.Conf) error {
			return rbacErr
		}

		if _, err := newServe(config, dependencies); !errors.Is(err, rbacErr) {
			t.Fatalf("newServe error = %v, want %v", err, rbacErr)
		}
	})

	t.Run("database", func(t *testing.T) {
		databaseErr := errors.New("database failed")
		dependencies := successfulServerDependencies(database, enforcer)
		dependencies.openDB = func(*conf.Mysql) (*gorm.DB, error) {
			return nil, databaseErr
		}

		if _, err := newServe(config, dependencies); !errors.Is(err, databaseErr) {
			t.Fatalf("newServe error = %v, want %v", err, databaseErr)
		}
	})

	t.Run("enforcer closes database", func(t *testing.T) {
		enforcerErr := errors.New("enforcer failed")
		closeErr := errors.New("close failed")
		closeCalls := 0
		dependencies := successfulServerDependencies(database, enforcer)
		dependencies.getAuth = func(*conf.Conf, *gorm.DB) (*casbin.Enforcer, error) {
			return nil, enforcerErr
		}
		dependencies.closeDB = func(*gorm.DB) error {
			closeCalls++
			return closeErr
		}

		_, err := newServe(config, dependencies)
		if !errors.Is(err, enforcerErr) || !errors.Is(err, closeErr) {
			t.Fatalf("newServe error = %v, want joined enforcer and close errors", err)
		}
		if closeCalls != 1 {
			t.Fatalf("close calls = %d, want 1", closeCalls)
		}
	})

	t.Run("migration closes database", func(t *testing.T) {
		migrationErr := errors.New("migration failed")
		closeErr := errors.New("close failed")
		closeCalls := 0
		dependencies := successfulServerDependencies(database, enforcer)
		dependencies.migrate = func(*WebServe) error {
			return migrationErr
		}
		dependencies.closeDB = func(*gorm.DB) error {
			closeCalls++
			return closeErr
		}

		_, err := newServe(config, dependencies)
		if !errors.Is(err, migrationErr) || !errors.Is(err, closeErr) {
			t.Fatalf("newServe error = %v, want joined migration and close errors", err)
		}
		if closeCalls != 1 {
			t.Fatalf("close calls = %d, want 1", closeCalls)
		}
	})
}

func TestNewServeInitializesSupportedLocales(t *testing.T) {
	for _, locale := range []string{"en", "zh", "unsupported"} {
		t.Run(locale, func(t *testing.T) {
			config := conf.NewConf()
			config.Locale = locale
			service, err := newServe(
				config,
				successfulServerDependencies(&gorm.DB{}, &casbin.Enforcer{}),
			)
			if err != nil {
				t.Fatalf("newServe error: %v", err)
			}
			if service.validate == nil {
				t.Fatal("validator was not initialized")
			}
		})
	}
}

func successfulServerDependencies(database *gorm.DB, enforcer *casbin.Enforcer) serverDependencies {
	return serverDependencies{
		ensureRbac: func(*conf.Conf) error {
			return nil
		},
		openDB: func(*conf.Mysql) (*gorm.DB, error) {
			return database, nil
		},
		getAuth: func(*conf.Conf, *gorm.DB) (*casbin.Enforcer, error) {
			return enforcer, nil
		},
		migrate: func(*WebServe) error {
			return nil
		},
		closeDB: func(*gorm.DB) error {
			return nil
		},
	}
}

func TestGormDbValidationAndConnectionFailures(t *testing.T) {
	if _, err := gormDbWithOpen(nil, nil); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("nil config error = %v, want %v", err, e.ErrConfigInvalid)
	}
	if _, err := gormDb(nil); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("gormDb nil config error = %v, want %v", err, e.ErrConfigInvalid)
	}
	if _, err := gormDbWithOpen(&conf.Mysql{}, nil); !errors.Is(err, e.ErrDbTableNameEmpty) {
		t.Fatalf("empty database error = %v, want %v", err, e.ErrDbTableNameEmpty)
	}

	mysqlConfig := conf.NewConf().Mysql
	openErr := errors.New("open failed")
	if _, err := gormDbWithOpen(mysqlConfig, func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
		return nil, openErr
	}); !errors.Is(err, openErr) {
		t.Fatalf("open error = %v, want %v", err, openErr)
	}

	if _, err := gormDbWithOpen(mysqlConfig, func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
		return &gorm.DB{Config: &gorm.Config{}}, nil
	}); err == nil {
		t.Fatal("invalid GORM database did not return an error")
	}

	pingErr := errors.New("ping failed")
	database := openTestGorm(t, &testSQLDriver{
		ping: func() error {
			return pingErr
		},
	})
	if _, err := gormDbWithOpen(mysqlConfig, func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
		return database, nil
	}); !errors.Is(err, pingErr) {
		t.Fatalf("ping error = %v, want %v", err, pingErr)
	}
}

func TestGormDbConfiguresConnectionPool(t *testing.T) {
	config := *conf.NewConf().Mysql
	config.MaxIdleConns = 2
	config.MaxOpenConns = 4
	database := openTestGorm(t, &testSQLDriver{})

	got, err := gormDbWithOpen(&config, func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
		return database, nil
	})
	if err != nil {
		t.Fatalf("gormDbWithOpen error: %v", err)
	}
	sqlDB, err := got.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	if sqlDB.Stats().MaxOpenConnections != config.MaxOpenConns {
		t.Fatalf("max open connections = %d, want %d", sqlDB.Stats().MaxOpenConnections, config.MaxOpenConns)
	}
}

func TestMysqlLogTargetExcludesCredentials(t *testing.T) {
	config := &conf.Mysql{
		Path:     "database.example:3306",
		DbName:   "service",
		Username: "private-user",
		Password: "private-password",
	}
	target := mysqlLogTarget(config)
	if !strings.Contains(target, config.Path) || !strings.Contains(target, config.DbName) {
		t.Fatalf("log target = %q, want address and database", target)
	}
	if strings.Contains(target, config.Username) || strings.Contains(target, config.Password) {
		t.Fatalf("log target contains credentials: %q", target)
	}
}

func TestCloseGormDB(t *testing.T) {
	if err := closeGormDB(nil); err != nil {
		t.Fatalf("close nil database: %v", err)
	}
	if err := closeGormDB(&gorm.DB{Config: &gorm.Config{}}); err == nil {
		t.Fatal("close invalid database did not return an error")
	}

	database := openTestGorm(t, &testSQLDriver{})
	if err := closeGormDB(database); err != nil {
		t.Fatalf("close database: %v", err)
	}
}

func TestRunEReturnsServingError(t *testing.T) {
	listenErr := errors.New("listen failed")
	server := &fakeServe{listenErr: listenErr}
	service := &WebServe{
		serve:  server,
		conf:   conf.NewConf(),
		engine: gin.New(),
	}

	if err := service.RunE(); !errors.Is(err, listenErr) {
		t.Fatalf("RunE error = %v, want %v", err, listenErr)
	}
	if server.listenCalls != 1 {
		t.Fatalf("ListenAndServe calls = %d, want 1", server.listenCalls)
	}
}

func TestRunLogsServingError(t *testing.T) {
	originalWriter := log.Writer()
	log.SetOutput(&strings.Builder{})
	t.Cleanup(func() {
		log.SetOutput(originalWriter)
	})

	server := &fakeServe{listenErr: errors.New("listen failed")}
	service := &WebServe{
		serve:  server,
		conf:   conf.NewConf(),
		engine: gin.New(),
	}
	service.Run()
	if server.listenCalls != 1 {
		t.Fatalf("ListenAndServe calls = %d, want 1", server.listenCalls)
	}
}

func TestRunEAcceptsServerClosed(t *testing.T) {
	service := &WebServe{
		serve:  &fakeServe{listenErr: http.ErrServerClosed},
		conf:   conf.NewConf(),
		engine: gin.New(),
	}
	if err := service.RunE(); err != nil {
		t.Fatalf("RunE error = %v, want nil", err)
	}
}

func TestRunECreatesAndStopsServer(t *testing.T) {
	config := conf.NewConf()
	config.System.Addr = "127.0.0.1:0"
	service := &WebServe{
		conf:   config,
		engine: gin.New(),
	}
	result := make(chan error, 1)
	go func() {
		result <- service.RunE()
	}()

	var server serve
	deadline := time.Now().Add(time.Second)
	for server == nil && time.Now().Before(deadline) {
		service.serveMu.RLock()
		server = service.serve
		service.serveMu.RUnlock()
		if server == nil {
			time.Sleep(time.Millisecond)
		}
	}
	if server == nil {
		t.Fatal("RunE did not initialize the HTTP server")
	}
	if err := service.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown error: %v", err)
	}
	<-result
}

func TestRunERequiresInitializedService(t *testing.T) {
	var service *WebServe
	if err := service.RunE(); !errors.Is(err, ErrServerNotInitialized) {
		t.Fatalf("nil RunE error = %v, want %v", err, ErrServerNotInitialized)
	}
	if err := (&WebServe{}).RunE(); !errors.Is(err, ErrServerNotInitialized) {
		t.Fatalf("empty RunE error = %v, want %v", err, ErrServerNotInitialized)
	}
}

func TestShutdownDelegatesAndValidatesContext(t *testing.T) {
	shutdownErr := errors.New("shutdown failed")
	server := &fakeServe{shutdownErr: shutdownErr}
	service := &WebServe{serve: server}

	if err := service.Shutdown(nil); !errors.Is(err, ErrShutdownContextRequired) {
		t.Fatalf("Shutdown(nil) error = %v, want %v", err, ErrShutdownContextRequired)
	}
	if err := service.Shutdown(context.Background()); !errors.Is(err, shutdownErr) {
		t.Fatalf("Shutdown error = %v, want %v", err, shutdownErr)
	}
	if server.shutdownCalls != 1 {
		t.Fatalf("Shutdown calls = %d, want 1", server.shutdownCalls)
	}

	var nilService *WebServe
	if err := nilService.Shutdown(context.Background()); !errors.Is(err, ErrServerNotInitialized) {
		t.Fatalf("nil service Shutdown error = %v, want %v", err, ErrServerNotInitialized)
	}
}

func TestShutdownClosesDatabaseOnce(t *testing.T) {
	closeErr := errors.New("close failed")
	closeCalls := 0
	database := openTestGorm(t, &testSQLDriver{
		close: func() error {
			closeCalls++
			return closeErr
		},
	})
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("open test connection: %v", err)
	}

	service := &WebServe{db: database}
	if err := service.Shutdown(context.Background()); !errors.Is(err, closeErr) {
		t.Fatalf("first Shutdown error = %v, want %v", err, closeErr)
	}
	if err := service.Shutdown(context.Background()); !errors.Is(err, closeErr) {
		t.Fatalf("second Shutdown error = %v, want %v", err, closeErr)
	}
	if closeCalls != 1 {
		t.Fatalf("database close calls = %d, want 1", closeCalls)
	}
}
