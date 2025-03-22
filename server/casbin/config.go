package casbin

import (
	"fmt"
	"path/filepath"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/g"
)

// Remove del config file
func Remove() error {
	casbinPath := getCasbinPath()
	if dir.IsExist(casbinPath) && dir.IsFile(casbinPath) {
		return dir.Remove(casbinPath)
	}
	return nil
}

func getCasbinPath() string {
	return filepath.Join(dir.GetCurrentAbPath(), g.CasbinFileName)
}

// init initialize config file
// - initialize casbin's config file as rbac_model.conf name
func init() {
	casbinPath := getCasbinPath()
	if !dir.IsExist(casbinPath) { // casbin rbac_model.conf file
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
		_, err := dir.WriteBytes(casbinPath, rbacModelConf)
		if err != nil {
			panic(fmt.Errorf("initialize casbin rbac_model.conf file return error: %w ", err))
		}
	}
}
