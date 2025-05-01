package admin

import (
	"path/filepath"
	"strconv"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/conf"
	"gorm.io/gorm"
)

// getEnforcer get casbin.Enforcer
func getEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	if db == nil {
		return nil, gorm.ErrInvalidDB
	}
	c, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule") // Your driver and data source.
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(filepath.Join(dir.GetCurrentAbPath(), conf.RbacName), c)
	if err != nil {
		return nil, err
	}

	if err = enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return enforcer, nil
}

// GetRolesForUser get user's roles
func (ws *WebServe) GetRolesForUser(uid uint) []string {
	uids, err := ws.Auth().GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		return []string{}
	}

	return uids
}

// ClearCasbin clean rules
func (ws *WebServe) ClearCasbin(v int, p ...string) error {
	_, err := ws.Auth().RemoveFilteredPolicy(v, p...)
	if err != nil {
		return err
	}
	return nil
}
