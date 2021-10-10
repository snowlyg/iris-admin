package v1

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/modules/v1/auth"
	"github.com/snowlyg/iris-admin/modules/v1/file"
	"github.com/snowlyg/iris-admin/modules/v1/oplog"
	"github.com/snowlyg/iris-admin/modules/v1/perm"
	"github.com/snowlyg/iris-admin/modules/v1/role"
	"github.com/snowlyg/iris-admin/modules/v1/user"
)

// Party v1 模块
func Party() func(v1 iris.Party) {
	return func(v1 iris.Party) {
		v1.PartyFunc("/users", user.Party())
		v1.PartyFunc("/roles", role.Party())
		v1.PartyFunc("/perms", perm.Party())
		v1.PartyFunc("/file", file.Party())
		v1.PartyFunc("/auth", auth.Party())
		v1.PartyFunc("/oplog", oplog.Party())
	}
}
