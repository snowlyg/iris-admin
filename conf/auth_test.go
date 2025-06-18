package conf

import (
	"testing"

	"github.com/snowlyg/helper/dir"
)

func TestNewRbacModel(t *testing.T) {
	conf := new(Conf)
	cfp := conf.casbinFilePath()
	if cfp == "" {
		t.Errorf("rbac model path:%s empty", cfp)
	}
	conf.newRbacModel()
	if !dir.IsExist(cfp) {
		t.Errorf("%s not exist after conf not init and new rbac model %t", cfp, dir.IsExist(cfp))
	}
	conf.RemoveRbacModel()
	if dir.IsExist(cfp) {
		t.Errorf("%s exist after conf not init and new rbac model %t", cfp, dir.IsExist(cfp))
	}
}
