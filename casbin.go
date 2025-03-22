package admin

import (
	"github.com/casbin/casbin/v2"
)

var (
	// once     sync.Once
	enforcer *casbin.Enforcer
)

// // Instance casbin instance
// func Instance() *casbin.Enforcer {
// 	once.Do(func() {
// 		enforcer = getEnforcer()
// 	})
// 	return enforcer
// }

// getEnforcer get casbin.Enforcer
func getEnforcer() *casbin.Enforcer {
	// if Instance() == nil {
	// 	return nil
	// }
	// c, err := gormadapter.NewAdapterByDBUseTableName(Instance(), "", "casbin_rule") // Your driver and data source.
	// if err != nil {
	// 	return nil
	// }

	// enforcer, err := casbin.NewEnforcer(filepath.Join(dir.GetCurrentAbPath(), admin.CasbinFileName), c)
	// if err != nil {
	// 	return nil
	// }

	// if enforcer == nil {
	// 	return nil
	// }

	// err = enforcer.LoadPolicy()
	// if err != nil {
	// 	return nil
	// }

	return enforcer
}

// // GetRolesForUser get user's roles
// func GetRolesForUser(uid uint) []string {
// 	uids, err := Instance().GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
// 	if err != nil {
// 		return []string{}
// 	}

// 	return uids
// }

// // ClearCasbin clean rules
// func ClearCasbin(v int, p ...string) error {
// 	_, err := Instance().RemoveFilteredPolicy(v, p...)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
