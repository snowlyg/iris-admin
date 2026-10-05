package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorWithStatusUsesHTTPStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ErrorWithStatus(http.StatusBadRequest, "invalid request", ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	response := decodeResponse(t, recorder)
	if response.Code != http.StatusBadRequest || response.Msg != "invalid request" {
		t.Fatalf("response = %#v, want status 400 and safe message", response)
	}
}

func TestErrorWithStatusRejectsNonErrorStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ErrorWithStatus(http.StatusOK, "invalid status", ctx)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	response := decodeResponse(t, recorder)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("response status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}

func TestErrorWithStatusAndDataPreservesDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ErrorWithStatusAndData(http.StatusUnprocessableEntity, map[string]string{"field": "required"}, "validation failed", ctx)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	response := decodeResponse(t, recorder)
	data, ok := response.Data.(map[string]any)
	if !ok || data["field"] != "required" {
		t.Fatalf("response data = %#v, want validation details", response.Data)
	}
}

func TestFailWithMessageRetainsLegacyHTTPStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	FailWithMessage("legacy error", ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want legacy status %d", recorder.Code, http.StatusOK)
	}
	response := decodeResponse(t, recorder)
	if response.Code != http.StatusBadRequest || response.Msg != "legacy error" {
		t.Fatalf("response = %#v, want body status 400 and legacy message", response)
	}
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) Response {
	t.Helper()
	var response Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return response
}
