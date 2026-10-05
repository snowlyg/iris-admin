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

func TestStart(t *testing.T) {
	t.Setenv("IRIS_ADMIN_WEB_ADDR", "127.0.0.1:18088")
	c, err := conf.LoadConfFile(filepath.Join(t.TempDir(), "iris_admin.json"))
	if err != nil {
		t.Fatal(err)
	}
	serve, err := NewServe(c)
	if err != nil {
		t.Fatal(err)
	}
	go serve.Run()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := serve.Shutdown(ctx); err != nil {
			t.Errorf("shutdown server: %v", err)
		}
	}()

	time.Sleep(3 * time.Second)

	resp, err := http.Get("http://127.0.0.1:18088")
	if err != nil {
		t.Fatalf("test web start get %v", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("test web start get %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("test web start want [%d] but get [%d]", http.StatusNotFound, resp.StatusCode)
	}
}
