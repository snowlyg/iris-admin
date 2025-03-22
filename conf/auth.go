package conf

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/snowlyg/helper/dir"
)

const rbacName = "rbac_model.conf"

// Remove del config file
func (conf *Conf) RemoveRbacModel() error {
	p := conf.getPath()
	if filepath.Base(p) != rbacName {
		return nil
	}
	if dir.IsExist(p) && dir.IsFile(p) {
		return dir.Remove(p)
	}
	return nil
}

// getPath
func (conf *Conf) getPath() string {
	return filepath.Join(dir.GetCurrentAbPath(), rbacName)
}

// newRbacModel initialize casbin's config file as rbac_model.conf name
func (conf *Conf) newRbacModel() {
	if dir.IsExist(conf.getPath()) {
		// casbin rbac_model.conf file
		log.Printf("rbac_model.conf file is existed.")
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
	if _, err := dir.WriteBytes(conf.getPath(), rbacModelConf); err != nil {
		panic(fmt.Errorf("initialize casbin rbac_model.conf file return error: %w ", err))
	}
}
