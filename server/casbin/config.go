package casbin

import (
	"fmt"
	"path/filepath"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/g"
)

// Remove 删除配置文件
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

// new 初始化系统配置
// - 第一次初始化系统配置，会自动生成casbin 的规则文件 rbac_model.conf
func new() {
	casbinPath := getCasbinPath()
	fmt.Printf("casbin rbac_model.conf 位于： %s\n\n", casbinPath)
	if !dir.IsExist(casbinPath) { // casbin rbac_model.conf 文件
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
			panic(fmt.Errorf("初始化 casbin rbac_model.conf 文件错误: %w ", err))
		}
	}
}
