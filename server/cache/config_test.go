package cache

import (
	"os"
	"testing"
)

func TestIsExist(t *testing.T) {
	t.Run("测试redis配置初始化方法", func(t *testing.T) {
		redisPwd := os.Getenv("redisPwd")
		CONFIG.Password = redisPwd
		Instance()
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
