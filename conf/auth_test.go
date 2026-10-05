package conf

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/e"
)

func TestNewRbacModel(t *testing.T) {
	conf := new(Conf)
	conf.configDir = t.TempDir()
	cfp := conf.casbinFilePath()
	if cfp == "" {
		t.Errorf("rbac model path:%s empty", cfp)
	}
	if err := conf.EnsureRbacModel(); err != nil {
		t.Fatalf("ensure rbac model: %v", err)
	}
	if !dir.IsExist(cfp) {
		t.Errorf("%s not exist after conf not init and new rbac model %t", cfp, dir.IsExist(cfp))
	}
	if err := conf.EnsureRbacModel(); err != nil {
		t.Fatalf("ensure existing rbac model: %v", err)
	}
	if err := conf.RemoveRbacModel(); err != nil {
		t.Fatalf("remove rbac model: %v", err)
	}
	if dir.IsExist(cfp) {
		t.Fatalf("%s still exists after removal", cfp)
	}
	if err := conf.RemoveRbacModel(); err != nil {
		t.Fatalf("remove missing rbac model: %v", err)
	}
}

func TestRbacModelRejectsInvalidConfiguration(t *testing.T) {
	var conf *Conf
	if err := conf.EnsureRbacModel(); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("EnsureRbacModel nil error = %v, want %v", err, e.ErrConfigInvalid)
	}
	if err := conf.RemoveRbacModel(); !errors.Is(err, e.ErrConfigInvalid) {
		t.Fatalf("RemoveRbacModel nil error = %v, want %v", err, e.ErrConfigInvalid)
	}
}

func TestEnsureRbacModelReportsWriteFailure(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(blocker, []byte("block"), 0o600); err != nil {
		t.Fatal(err)
	}
	conf := newDefaultConf()
	conf.configDir = blocker
	if err := conf.EnsureRbacModel(); err == nil {
		t.Fatal("EnsureRbacModel ignored write failure")
	}
}
