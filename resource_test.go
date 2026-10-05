package admin

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/snowlyg/iris-admin/conf"
)

func baseServer(t *testing.T) (*WebServe, error) {
	t.Helper()
	c, err := conf.LoadConfFile(filepath.Join(t.TempDir(), "iris_admin.json"))
	if err != nil {
		return nil, err
	}
	c.System.Addr = "127.0.0.1:8088"
	s, err := NewServe(c)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func TestWebServeResource(t *testing.T) {
	s, err := baseServer(t)
	if err != nil {
		t.Fatalf("baseServer err:%s\n", err.Error())
	}

	engine := s.Engine()
	// add group api v1
	v1 := engine.Group("/api/v1")
	{
		s.Resource(v1, new(Router))
	}
	go s.Run()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			t.Errorf("shutdown server: %v", err)
		}
	}()

	time.Sleep(3 * time.Second)

	resp, err := http.Get("http://127.0.0.1:8088/api/v1/routers/list")
	if err != nil {
		t.Fatalf("get err:%s\n", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("test web start want [%d] but get [%d]", http.StatusOK, resp.StatusCode)
	}
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("read body err:%s\n", err.Error())
	}
	t.Logf("result:%s\n", string(result))
}
