package v1

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/rate"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/modules/v1/auth"
	"github.com/snowlyg/iris-admin/modules/v1/initdb"
	"github.com/snowlyg/iris-admin/modules/v1/perm"
	"github.com/snowlyg/iris-admin/modules/v1/role"
	"github.com/snowlyg/iris-admin/modules/v1/user"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party
func Party() module.WebModule {
	handler := func(v1 iris.Party) {
		if !g.CONFIG.Limit.Disable {
			limitV1 := rate.Limit(g.CONFIG.Limit.Limit, g.CONFIG.Limit.Burst, rate.PurgeEvery(time.Minute, 5*time.Minute))
			v1.Use(limitV1)
		}
	}
	modules := []module.WebModule{
		auth.Party(),
		initdb.Party(),
		role.Party(),
		perm.Party(),
		user.Party(),
	}
	return module.NewModule("/api/v1", handler, modules...)
}
