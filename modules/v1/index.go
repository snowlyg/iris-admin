package v1

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/rate"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/modules/initdb"
	"github.com/snowlyg/iris-admin/modules/perm"
	"github.com/snowlyg/iris-admin/modules/role"
	"github.com/snowlyg/iris-admin/modules/user"
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
		initdb.Party(),
		role.Party(),
		perm.Party(),
		user.Party(),
	}
	return module.NewModule("/api/v1", handler, modules...)
}
