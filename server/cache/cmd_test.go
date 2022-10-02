package cache

import (
	"testing"

	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

func TestInitConfig(t *testing.T) {
	defer zap_server.Remove()
	t.Run("test redis's file initialize", func(t *testing.T) {
		CONFIG.Password = g.TestRedisPwd
		err := InitConfig()
		if err != nil {
			t.Error(err)
			return
		}
		if !IsExist() {
			t.Errorf("config's files is not exist.")
			return
		}
		if err := Remove(); err != nil {
			t.Error(err)
		}
	})
}
