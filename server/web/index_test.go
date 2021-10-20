package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/snowlyg/iris-admin/server/web/web_iris"
)

func TestStart(t *testing.T) {
	addr := "127.0.0.1:8086"
	web_iris.CONFIG.System.Addr = addr
	wi := web_iris.Init()
	go func() {
		Start(wi)
	}()

	time.Sleep(3 * time.Second)

	t.Run("test web start", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://%s", addr))
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
