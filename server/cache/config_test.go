package cache

import (
	"os"
	"testing"

	"github.com/snowlyg/iris-admin/server/viper_server"
)

func TestIsExist(t *testing.T) {
	viper_server.Init(getViperConfig())
	t.Run("测试redis配置初始化方法", func(t *testing.T) {
		redisPwd := os.Getenv("redisPwd")
		CONFIG.Password = redisPwd
		if !IsExist() {
			t.Errorf("config's files is not exist.")
		}
	})
	t.Run("Test Remove function", func(t *testing.T) {
		if err := Remove(); err != nil {
			t.Error(err)
		}
		if IsExist() {
			t.Errorf("config's files remove is fail.")
		}
	})
}
