package web_iris

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

func TestRun(t *testing.T) {
	defer zap_server.Remove()
	defer web.Remove()
	web.CONFIG.System.Addr = "127.0.0.1:18085"
	go func() {
		web.Start(Init())
	}()

	time.Sleep(3 * time.Second)

	t.Run("test web run", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:18085/v0/version")
		if err != nil {
			t.Errorf("test web start get %v", err)
		}
		defer resp.Body.Close()
		s, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("test web start get %v", err)
		}
		if string(s) != "IRIS-ADMIN is running!!!" {
			t.Errorf("test web start want %s but get %s", "Not Found", string(s))
		}
	})
}
