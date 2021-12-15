package web_gin

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
)

func TestStart(t *testing.T) {
	web_iris.CONFIG.System.CacheType = "local"
	go func() {
		web.Start(Init())
	}()

	time.Sleep(3 * time.Second)

	t.Run("test web start", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:8088")
		if err != nil {
			t.Errorf("test web start get %v", err)
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("test web start get %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("test web start want [%d] but get [%d]", http.StatusNotFound, resp.StatusCode)
		}
	})
}
