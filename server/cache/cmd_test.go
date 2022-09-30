package cache

import (
	"os"
	"testing"
)

func TestInitConfig(t *testing.T) {
	t.Run("test redis's file initialize", func(t *testing.T) {
		redisPwd := os.Getenv("redisPwd")
		CONFIG.Password = redisPwd
		err := InitConfig()
		if err != nil {
			t.Error(err)
		}
		if !IsExist() {
			t.Errorf("config's files is not exist.")
		}
		if err := Remove(); err != nil {
			t.Error(err)
		}
	})
}
