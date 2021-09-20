package casbin

import (
	"path/filepath"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"go.uber.org/zap"
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

// GetEnforcer 获取 casbin.Enforcer
func GetEnforcer() *casbin.Enforcer {
	if database.Instance() == nil {
		g.ZAPLOG.Error("数据库未初始化")
		return nil
	}
	c, err := gormadapter.NewAdapterByDBUseTableName(database.Instance(), "", "casbin_rule") // Your driver and data source.
	if err != nil {
		g.ZAPLOG.Error("驱动初始化错误", zap.String("gormadapter.NewAdapterByDBUseTableName()", err.Error()))
		return nil
	}

	enforcer, err := casbin.NewEnforcer(filepath.Join(dir.GetCurrentAbPath(), g.CasbinFileName), c)
	if err != nil {
		g.ZAPLOG.Error("初始化失败", zap.String("casbin.NewEnforcer()", err.Error()))
		return nil
	}

	if enforcer == nil {
		g.ZAPLOG.Error("Casbin 未初始化")
		return nil
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		g.ZAPLOG.Error("加载规则失败", zap.String("casbin.LoadPolicy()", err.Error()))
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
