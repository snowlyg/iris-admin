package user

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"go.uber.org/zap"
)

type Request struct {
	BaseUser
	Password string `json:"password"`
	RoleIds  []uint `json:"role_ids"`
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
