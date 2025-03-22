package admin

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	// defer Remove()
	// CONFIG.System.Addr = "127.0.0.1:18088"
	go func() {
		Start(NewServe())
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
