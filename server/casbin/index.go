package casbin

import (
	"strconv"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
)

var (
	once     sync.Once
	enforcer *casbin.Enforcer
)

// Instance casbin 单例
func Instance() *casbin.Enforcer {
	once.Do(func() {
		enforcer = GetEnforcer()
	})
	return enforcer
}

// GetEnforcer
func GetEnforcer() *casbin.Enforcer {
	c, err := gormadapter.NewAdapterByDBUseTableName(database.Instance(), "", "casbin_rule") // Your driver and data source.
	if err != nil {
		return nil
	}

	enforcer, err := casbin.NewEnforcer(g.CasbinFileName, c)
	if err != nil {
		return nil
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		return nil
	}

	return enforcer
}

// GetRolesForUser 获取角色
func GetRolesForUser(uid uint) []string {
	uids, err := Instance().GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		return []string{}
	}

	return uids
}

// GetPermissionsForUser 获取角色权限
func GetPermissionsForUser(id string) [][]string {
	return Instance().GetPermissionsForUser(id)
}
