package admin

import (
	"database/sql/driver"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/conf"
)

func TestClassifyRoutes(t *testing.T) {
	routes := gin.RoutesInfo{
		{Method: http.MethodGet, Path: "/api/v1/users", Handler: "getUsers"},
		{Method: http.MethodPost, Path: "/api/v1/users", Handler: "createUser"},
		{Method: http.MethodPatch, Path: "/api/v1/users/:id", Handler: "patchUser"},
		{Method: http.MethodGet, Path: "/assets/*filepath", Handler: "static"},
	}

	permRoutes, otherRoutes := classifyRoutes(routes, " get ", "api/v1/users")
	if len(permRoutes) != 1 {
		t.Fatalf("permission routes = %d, want 1", len(permRoutes))
	}
	if key := newRouteKey(permRoutes[0].Path, permRoutes[0].Method); key != (routeKey{path: "/api/v1/users", method: http.MethodPost}) {
		t.Fatalf("permission route = %+v", permRoutes[0])
	}
	if len(otherRoutes) != 2 {
		t.Fatalf("other routes = %d, want 2", len(otherRoutes))
	}
	if key := newRouteKey(otherRoutes[0].Path, otherRoutes[0].Method); key != (routeKey{path: "/api/v1/users", method: http.MethodGet}) {
		t.Fatalf("excluded route = %+v", otherRoutes[0])
	}
	if key := newRouteKey(otherRoutes[1].Path, otherRoutes[1].Method); key != (routeKey{path: "/api/v1/users/:id", method: http.MethodPatch}) {
		t.Fatalf("unsupported method route = %+v", otherRoutes[1])
	}
}

func TestRouteExceptionsIgnoreInvalidPairs(t *testing.T) {
	if got := routeExceptions("GET;POST", "/one"); len(got) != 0 {
		t.Fatalf("mismatched route exceptions = %v, want empty", got)
	}
	if got := routeExceptions("", ""); len(got) != 0 {
		t.Fatalf("empty route exceptions = %v, want empty", got)
	}
}

func TestDiffRoutesUsesPathAndMethod(t *testing.T) {
	existingGet := &Router{Path: "/same", Method: "get"}
	existingGet.ID = 1
	existingOld := &Router{Path: "/old", Method: http.MethodPost}
	existingOld.ID = 2

	desired := []*Router{
		{Path: "/same", Method: http.MethodGet},
		{Path: "/same", Method: http.MethodPost},
	}
	deleteIDs, additions := diffRoutes([]*Router{existingGet, existingOld}, desired)

	if len(deleteIDs) != 1 || deleteIDs[0] != existingOld.ID {
		t.Fatalf("delete IDs = %v, want [%d]", deleteIDs, existingOld.ID)
	}
	if len(additions) != 1 {
		t.Fatalf("additions = %d, want 1", len(additions))
	}
	if key := newRouteKey(additions[0].Path, additions[0].Method); key != (routeKey{path: "/same", method: http.MethodPost}) {
		t.Fatalf("addition = %+v", additions[0])
	}
}

func TestDiffRoutesDeletesExistingWhenDesiredIsEmpty(t *testing.T) {
	existing := &Router{Path: "/old", Method: http.MethodGet}
	existing.ID = 9

	deleteIDs, additions := diffRoutes([]*Router{existing}, nil)
	if len(deleteIDs) != 1 || deleteIDs[0] != existing.ID {
		t.Fatalf("delete IDs = %v, want [%d]", deleteIDs, existing.ID)
	}
	if len(additions) != 0 {
		t.Fatalf("additions = %d, want 0", len(additions))
	}
}

func TestGroupRoutersHandlesDeleteFailure(t *testing.T) {
	execErr := errors.New("delete failed")
	var executed string
	db := openTestGorm(t, &testSQLDriver{
		query: func(string) (driver.Rows, error) {
			return routerRows([][]driver.Value{routerRow(9, "/old", http.MethodGet)}), nil
		},
		exec: func(query string) error {
			executed = query
			return execErr
		},
	})
	service := &WebServe{
		conf:   conf.NewConf(),
		db:     db,
		engine: gin.New(),
	}

	service.groupRouters()

	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(executed)), "UPDATE") {
		t.Fatalf("executed query = %q, want soft-delete UPDATE", executed)
	}
}

func TestGroupRoutersWithoutDatabaseStillClassifiesRoutes(t *testing.T) {
	engine := gin.New()
	engine.GET("/classified", func(*gin.Context) {})
	service := &WebServe{
		conf:   conf.NewConf(),
		engine: engine,
	}

	service.groupRouters()

	if len(service.permRoutes) != 1 || service.permRoutes[0].Path != "/classified" {
		t.Fatalf("permission routes = %+v, want classified route", service.permRoutes)
	}
}

func TestGroupRoutersHandlesCreateFailure(t *testing.T) {
	execErr := errors.New("create failed")
	var executed string
	db := openTestGorm(t, &testSQLDriver{
		query: func(string) (driver.Rows, error) {
			return routerRows(nil), nil
		},
		exec: func(query string) error {
			executed = query
			return execErr
		},
	})
	engine := gin.New()
	engine.GET("/new", func(*gin.Context) {})
	service := &WebServe{
		conf:   conf.NewConf(),
		db:     db,
		engine: engine,
	}

	service.groupRouters()

	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(executed)), "INSERT") {
		t.Fatalf("executed query = %q, want INSERT", executed)
	}
}

func TestGroupRoutersPersistsRouteChanges(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		var executed string
		db := openTestGorm(t, &testSQLDriver{
			query: func(string) (driver.Rows, error) {
				return routerRows([][]driver.Value{routerRow(9, "/old", http.MethodGet)}), nil
			},
			exec: func(query string) error {
				executed = query
				return nil
			},
		})
		service := &WebServe{conf: conf.NewConf(), db: db, engine: gin.New()}

		service.groupRouters()

		if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(executed)), "UPDATE") {
			t.Fatalf("executed query = %q, want soft-delete UPDATE", executed)
		}
	})

	t.Run("create", func(t *testing.T) {
		var executed string
		db := openTestGorm(t, &testSQLDriver{
			query: func(string) (driver.Rows, error) {
				return routerRows(nil), nil
			},
			exec: func(query string) error {
				executed = query
				return nil
			},
		})
		engine := gin.New()
		engine.GET("/new", func(*gin.Context) {})
		service := &WebServe{conf: conf.NewConf(), db: db, engine: engine}

		service.groupRouters()

		if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(executed)), "INSERT") {
			t.Fatalf("executed query = %q, want INSERT", executed)
		}
	})
}

func routerRows(values [][]driver.Value) driver.Rows {
	return &testSQLRows{
		columns: []string{"id", "created_at", "updated_at", "deleted_at", "path", "title", "group", "method"},
		values:  values,
	}
}

func routerRow(id int64, path, method string) []driver.Value {
	now := time.Now()
	return []driver.Value{id, now, now, nil, path, path, "", method}
}
