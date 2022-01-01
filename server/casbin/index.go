package casbin

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

var (
	once     sync.Once
	enforcer *casbin.Enforcer
)

// Instance casbin 单例
func Instance() *casbin.Enforcer {
	once.Do(func() {
		new()
		enforcer = GetEnforcer()
	})
	return enforcer
}

// GetEnforcer 获取 casbin.Enforcer
func GetEnforcer() *casbin.Enforcer {
	if database.Instance() == nil {
		zap_server.ZAPLOG.Error(database.ErrDatabaseInit.Error())
		return nil
	}
	c, err := gormadapter.NewAdapterByDBUseTableName(database.Instance(), "", "casbin_rule") // Your driver and data source.
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return nil
	}

	enforcer, err := casbin.NewEnforcer(filepath.Join(dir.GetCurrentAbPath(), g.CasbinFileName), c)
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return nil
	}

	if enforcer == nil {
		zap_server.ZAPLOG.Error("Casbin 未初始化")
		return nil
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
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

// ClearCasbin 清除权限
func ClearCasbin(v int, p ...string) error {
	_, err := Instance().RemoveFilteredPolicy(v, p...)
	if err != nil {
		return fmt.Errorf("清除权限失败 %w", err)
	}
	return nil
}
