package authority

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
)

// type Request struct {
// 	BaseAuthority
// 	Perms [][]string `json:"perms"`
// }

// func (req *Request) Request(ctx *gin.Context) error {
// 	if err := ctx.ShouldBindJSON(req); err != nil {
// 		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
// 		return orm.ErrParamValidate
// 	}
// 	return nil
// }

type ReqPaginate struct {
	orm.Paginate
	Name string `json:"name"`
}

type AuthorityCopyResponse struct {
	Authority      Authority `json:"authority"`
	OldAuthorityId uint      `json:"oldAuthorityId" binding:"required"`
}

func (req *AuthorityCopyResponse) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return orm.ErrParamValidate
	}
	return nil
}
