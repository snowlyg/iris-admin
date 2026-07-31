package conf

import (
	"fmt"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/e"
	"gorm.io/gorm"
)

const CasbinName = "rbac_model.conf"

// Remove del config file
func (conf *Conf) RemoveRbacModel() error {
	if conf == nil {
		return e.ErrConfigInvalid
	}
	p := conf.casbinFilePath()
	if filepath.Base(p) != CasbinName {
		return nil
	}
	if dir.IsExist(p) && dir.IsFile(p) {
		return dir.Remove(p)
	}
	return nil
}

func (conf *Conf) casbinFilePath() string {
	return filepath.Join(conf.configDirectory(), CasbinName)
}

// EnsureRbacModel initializes Casbin's model file when it is missing.
func (conf *Conf) EnsureRbacModel() error {
	if conf == nil {
		return e.ErrConfigInvalid
	}
	if dir.IsExist(conf.casbinFilePath()) {
		return nil
	}

	var rbacModelConf = []byte(`[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")`)
	if _, err := dir.WriteBytes(conf.casbinFilePath(), rbacModelConf); err != nil {
		return fmt.Errorf("initialize casbin rbac model: %w", err)
	}
	return nil
}

// getEnforcer get casbin.Enforcer
func (conf *Conf) GetEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	c, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule") // Your driver and data source.
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(conf.casbinFilePath(), c)
	if err != nil {
		return nil, err
	}
	if err = enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	return enforcer, nil
}
