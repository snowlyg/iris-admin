package cache

import "testing"

func TestIsExist(t *testing.T) {
	t.Run("测试redis配置初始化方法", func(t *testing.T) {
		err := InitConfig()
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Test IsExist function", func(t *testing.T) {
		if !IsExist() {
			t.Errorf("config's files is not exist.")
		}
	})
	t.Run("Test Remove function", func(t *testing.T) {
		if err := Remove(); err != nil {
			t.Error(err)
		}
	})
	t.Run("Test IsExist function again", func(t *testing.T) {
		if IsExist() {
			t.Errorf("config's files remove is fail.")
		}
	})
}
