package casbin

import (
	"testing"

	"github.com/snowlyg/helper/dir"
)

func TestRemove(t *testing.T) {
	t.Run("Test casbin instance", func(t *testing.T) {
		casbin := Instance()
		if casbin == nil {
			t.Error("casbin instance is nil")
		}
	})
	casbinPath := getCasbinPath()
	if !dir.IsExist(casbinPath) || !dir.IsFile(casbinPath) {
		t.Error("casbin file is not exist")
	}
	t.Run("Test casbin config remove", func(t *testing.T) {
		err := Remove()
		if err != nil {
			t.Error(err)
		}
	})
	if dir.IsExist(casbinPath) && dir.IsFile(casbinPath) {
		t.Error("casbin file is delete fail")
	}
}
