package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/web/web_iris/validate"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Model
type Model struct {
	Id        uint   `json:"id"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

// ReqId 获取id请求参数
type ReqId struct {
	Id uint `json:"id" param:"id"`
}

func (req *ReqId) Request(ctx interface{}) error {
	if c, ok := ctx.(iris.Context); ok {
		return req.irisReadParams(c)
	} else if c, ok := ctx.(*gin.Context); ok {
		return req.ginShouldBindJSON(c)
	} else {
		ctxTypeName := reflect.TypeOf(ctx).Name()
		return fmt.Errorf("Context [%s] 类型错误", ctxTypeName)
	}
}

func (req *ReqId) irisReadParams(ctx iris.Context) error {
	if err := ctx.ReadParams(req); err != nil {
		zap_server.ZAPLOG.Error("id参数获取失败", zap.String("ReadParams()", err.Error()))
		return ErrParamValidate
	}
	return nil
}

func (req *ReqId) ginShouldBindJSON(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap_server.ZAPLOG.Error("id参数获取失败", zap.String("ShouldBindJSON()", err.Error()))
		return ErrParamValidate
	}
	return nil
}

// Paginate 验证请求参数
type Paginate struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	OrderBy  string `json:"orderBy" form:"orderBy"`
	Sort     string `json:"sort" form:"sort"`
}

func (req *Paginate) Request(ctx interface{}) error {
	if c, ok := ctx.(iris.Context); ok {
		return req.irisReadQuerys(c)
	} else if c, ok := ctx.(*gin.Context); ok {
		return req.ginShouldBind(c)
	} else {
		ctxTypeName := reflect.TypeOf(ctx).Name()
		return fmt.Errorf("Context [%s] 类型错误", ctxTypeName)
	}
}

func (req *Paginate) irisReadQuerys(ctx iris.Context) error {
	if err := ctx.ReadQuery(req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			zap_server.ZAPLOG.Error("参数验证失败", zap.String("ValidRequest()", strings.Join(errs, ";")))
			return ErrParamValidate
		}
	}
	return nil
}

func (req *Paginate) ginShouldBind(ctx *gin.Context) error {
	if err := ctx.ShouldBind(req); err != nil {
		zap_server.ZAPLOG.Error("id参数获取失败", zap.String("ShouldBind()", err.Error()))
		return ErrParamValidate
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

	ErrParamValidate = errors.New("参数验证失败")
	ErrPaginateParam = errors.New("分页查询参数缺失")
)
