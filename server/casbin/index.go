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
	"github.com/snowlyg/iris-admin/server/zap_server"
)

var (
	once     sync.Once
	enforcer *casbin.Enforcer
)

// Instance casbin instance
func Instance() *casbin.Enforcer {
	once.Do(func() {
		enforcer = getEnforcer()
	})
	return enforcer
}

// getEnforcer get casbin.Enforcer
func getEnforcer() *casbin.Enforcer {
	if database.Instance() == nil {
		zap_server.ZAPLOG.Error(database.ErrDatabaseInit.Error())
		return nil
	}
	c, err := gormadapter.NewAdapterByDBUseTableName(database.Instance(), "", "casbin_rule") // Your driver and data source.
	if err != nil {
		return nil
	}

	enforcer, err := casbin.NewEnforcer(filepath.Join(dir.GetCurrentAbPath(), g.CasbinFileName), c)
	if err != nil {
		return nil
	}

	if enforcer == nil {
		zap_server.ZAPLOG.Error("Casbin init")
		return nil
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		return nil
	}

	return enforcer
}

// GetRolesForUser get user's roles
func GetRolesForUser(uid uint) []string {
	uids, err := Instance().GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		return []string{}
	}

	return uids
}

// ClearCasbin clean rules
func ClearCasbin(v int, p ...string) error {
	_, err := Instance().RemoveFilteredPolicy(v, p...)
	if err != nil {
		return err
	}
	return nil
}
