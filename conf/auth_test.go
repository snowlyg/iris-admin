package conf

import (
	"testing"

	"github.com/snowlyg/helper/dir"
)

func TestNewRbacModel(t *testing.T) {
	conf := new(Conf)
	if conf.getPath() == "" {
		t.Errorf("rbac model path:%s empty", conf.getPath())
	}
	conf.newRbacModel()
	if !dir.IsExist(RbacName) {
		t.Error("rbac_model.conf not exist after conf not init and new rbac model")
	}
	conf.RemoveRbacModel()
	if dir.IsExist(RbacName) {
		t.Error("rbac_model.conf exist after conf not init and new rbac model")
	}
}
