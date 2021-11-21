package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
)

type Request struct {
	BaseAdmin
	Password     string `json:"password"`
	AuthorityIds []uint `json:"authorityIds" binding:"required"`
}

func (req *Request) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return orm.ErrParamValidate
	}
	return nil
}

type ReqPaginate struct {
	orm.Paginate
	Name string `json:"name"`
}
