package admin

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/snowlyg/iris-admin/conf"
)

func TestStart(t *testing.T) {
	go func() {
		os.Setenv("IRIS_ADMIN_WEB_ADDR", "127.0.0.1:18088")
		c := conf.NewConf()
		if serve, err := NewServe(c); err != nil {
			t.Error(err.Error())
		} else {
			if err := serve.InitRouter(); err != nil {
				t.Error(err.Error())
			}
			serve.Run()
		}
	}()

	time.Sleep(3 * time.Second)

	resp, err := http.Get("http://127.0.0.1:18088")
	if err != nil {
		t.Errorf("test web start get %v", err)
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
