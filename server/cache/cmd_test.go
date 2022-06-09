package cache

import "testing"

func TestInitConfig(t *testing.T) {
	t.Run("测试redis配置初始化方法", func(t *testing.T) {
		err := InitConfig()
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("测试 initConfig()", func(t *testing.T) {
		err := initConfig()
		if err != nil {
			t.Error(err)
		}
	})
}
