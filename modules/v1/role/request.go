package role

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/casbin"
	"go.uber.org/zap"
)

type Request struct {
	BaseRole
	Perms casbin.PermsCollection `json:"perms"`
}

func (req *Request) Request(ctx iris.Context) error {
	if err := ctx.ReadJSON(req); err != nil {
		g.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return g.ErrParamValidate
	}
	return nil
}

type ReqPaginate struct {
	g.Paginate
	Name string `json:"name"`
}
