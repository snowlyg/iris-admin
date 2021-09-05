package casbin

import (
	"path/filepath"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gookit/color"
	"github.com/snowlyg/helper/dir"
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
	if database.Instance() == nil {
		color.Danger.Println("数据库初始化为空")
		return nil
	}
	c, err := gormadapter.NewAdapterByDBUseTableName(database.Instance(), "", "casbin_rule") // Your driver and data source.
	if err != nil {
		color.Danger.Printf("Casbin 驱动初始化错误 %v \n", err)
		return nil
	}

	enforcer, err := casbin.NewEnforcer(filepath.Join(dir.GetCurrentAbPath(), g.CasbinFileName), c)
	if err != nil {
		color.Danger.Printf("Casbin 初始化失败 %v\n", err)
		return nil
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		color.Danger.Printf("Casbin 加载规则失败 %\n", err)
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
