package web_iris

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/snowlyg/iris-admin/server/web"
)

func TestStart(t *testing.T) {
	go func() {
		web.Start(Init())
	}()

	time.Sleep(3 * time.Second)

	t.Run("test web start", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:8085")
		if err != nil {
			t.Errorf("test web start get %v", err)
		}
		defer resp.Body.Close()
		s, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("test web start get %v", err)
		}
		if string(s) != "Not Found" {
			t.Errorf("test web start want %s but get %s", "Not Found", string(s))
		}
	})
}
