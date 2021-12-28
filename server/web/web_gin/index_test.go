package web_gin

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/snowlyg/iris-admin/server/web"
)

func TestStart(t *testing.T) {
	defer web.Remove()
	web.CONFIG.System.Addr = "localhost:18088"
	go func() {
		web.Start(Init())
	}()

	time.Sleep(3 * time.Second)

	t.Run("test web start", func(t *testing.T) {
		resp, err := http.Get("http://localhost:18088")
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
