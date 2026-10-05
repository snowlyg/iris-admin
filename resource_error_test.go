package admin

import (
	"database/sql/driver"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResourceWithoutDatabaseReturnsSafeServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	service := &WebServe{}
	service.Resource(engine.Group("/api/v1"), new(Router))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/routers/list", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	response := decodeResponse(t, recorder)
	if response.Code != http.StatusInternalServerError || response.Msg != databaseUnavailableMessage {
		t.Fatalf("response = %#v, want safe database-unavailable error", response)
	}
}

func TestResourceDatabaseQueryFailureReturnsSafeServerError(t *testing.T) {
	queryErr := errors.New("private database query details")
	db := openTestGorm(t, &testSQLDriver{
		query: func(string) (driver.Rows, error) {
			return nil, queryErr
		},
	})
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	service := &WebServe{db: db}
	service.Resource(engine.Group("/api/v1"), new(Router))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/routers/list", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	response := decodeResponse(t, recorder)
	if response.Msg != databaseQueryFailedMessage {
		t.Fatalf("response message = %q, want %q", response.Msg, databaseQueryFailedMessage)
	}
	if strings.Contains(recorder.Body.String(), queryErr.Error()) {
		t.Fatal("response leaked the database error")
	}
}
