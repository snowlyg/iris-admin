package g

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/web/validate"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Model
type Model struct {
	Id        uint   `json:"id"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

// ReqId 验证请求参数
type ReqId struct {
	Id uint `json:"id" param:"id"`
}

func (req *ReqId) Request(ctx iris.Context) error {
	if err := ctx.ReadParams(req); err != nil {
		ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return orm.ErrParamValidate
	}
	return nil
}

// Paginate 验证请求参数
type Paginate struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	OrderBy  string `json:"orderBy"`
	Sort     string `json:"sort"`
}

func (req *Paginate) Request(ctx iris.Context) error {
	if err := ctx.ReadQuery(req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			ZAPLOG.Error("参数验证失败", zap.String("ValidRequest()", strings.Join(errs, ";")))
			return orm.ErrParamValidate
		}
	}
	return nil
}

// PaginateScope 分页 scope
func (req *Paginate) PaginateScope() func(db *gorm.DB) *gorm.DB {
	return scope.PaginateScope(req.Page, req.PageSize, req.Sort, req.OrderBy)
}

// Response
type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

// ErrMsg
type ErrMsg struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

var (
	NoErr         = ErrMsg{2000, "请求成功"}
	NeedInitErr   = ErrMsg{2001, "前往初始化数据库"}
	AuthErr       = ErrMsg{4001, "认证错误"}
	AuthExpireErr = ErrMsg{4002, "token 过期，请刷新token"}
	AuthActionErr = ErrMsg{4003, "权限错误"}
	ParamErr      = ErrMsg{4004, "参数解析失败，请联系管理员"}
	SystemErr     = ErrMsg{5000, "系统错误，请联系管理员"}
	DataEmptyErr  = ErrMsg{5001, "数据为空"}
	TokenCacheErr = ErrMsg{5002, "TOKEN CACHE 错误"}
)
