package conf

import (
	"fmt"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/snowlyg/helper/dir"
	"gorm.io/gorm"
)

const CasbinName = "rbac_model.conf"

// Remove del config file
func (conf *Conf) RemoveRbacModel() error {
	p := conf.casbinFilePath()
	if filepath.Base(p) != CasbinName {
		return nil
	}
	if dir.IsExist(p) && dir.IsFile(p) {
		return dir.Remove(p)
	}
	return nil
}

// casbinFilePath
func (conf *Conf) casbinFilePath() string {
	return filepath.Join(dir.GetCurrentAbPath(), ConfigDir, CasbinName)
}

// newRbacModel initialize casbin's config file as rbac_model.conf name
func (conf *Conf) newRbacModel() {
	if dir.IsExist(conf.casbinFilePath()) {
		// casbin rbac_model.conf file
		// log.Printf("rbac_model.conf file is existed.")
		return
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
		panic(fmt.Errorf("initialize casbin rbac_model.conf file return error: %w ", err))
	}
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
