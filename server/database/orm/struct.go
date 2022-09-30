package orm

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/web/web_iris/validate"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"gorm.io/gorm"
)

type ErrMsg struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

var (
	ErrParamValidate      = errors.New("参数验证失败")
	ErrPaginateParam      = errors.New("分页查询参数缺失")
	ErrUnSupportFramework = errors.New("不支持的框架")
)

// Model
type Model struct {
	Id        uint   `json:"id" uri:"id" form:"id" param:"id"`
	UpdatedAt string `json:"updatedAt" uri:"updatedAt" form:"updatedAt" param:"updatedAt"`
	CreatedAt string `json:"createdAt" uri:"createdAt" form:"createdAt" param:"createdAt"`
	DeletedAt string `json:"deletedAt" uri:"deletedAt" form:"deletedAt" param:"deletedAt"`
}

// ReqId the struct has used to get id form the context of every query
type ReqId struct {
	Id uint `json:"id" uri:"id" form:"id" param:"id"`
}

// Request get id data form the context of every query
func (req *ReqId) Request(ctx interface{}) error {
	if c, ok := ctx.(iris.Context); ok {
		return req.irisReadParams(c)
	} else if c, ok := ctx.(*gin.Context); ok {
		return req.ginShouldBindUri(c)
	} else {
		return ErrUnSupportFramework
	}
}

func (req *ReqId) irisReadParams(ctx iris.Context) error {
	if err := ctx.ReadParams(req); err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return ErrParamValidate
	}
	return nil
}

func (req *ReqId) ginShouldBindUri(ctx *gin.Context) error {
	if err := ctx.ShouldBindUri(req); err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return ErrParamValidate
	}
	return nil
}

// Paginate param for paginate query
type Paginate struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	OrderBy  string `json:"orderBy" form:"orderBy"`
	Sort     string `json:"sort" form:"sort"`
}

// Request
func (req *Paginate) Request(ctx interface{}) error {
	if c, ok := ctx.(iris.Context); ok {
		return req.irisReadQuerys(c)
	} else if c, ok := ctx.(*gin.Context); ok {
		return req.ginShouldBind(c)
	} else {
		return ErrUnSupportFramework
	}
}

func (req *Paginate) irisReadQuerys(ctx iris.Context) error {
	if err := ctx.ReadQuery(req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			zap_server.ZAPLOG.Error(strings.Join(errs, ";"))
			return ErrParamValidate
		}
	}
	return nil
}

func (req *Paginate) ginShouldBind(ctx *gin.Context) error {
	if err := ctx.ShouldBind(req); err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		return ErrParamValidate
	}
	return nil
}

// PaginateScope paginate scope
func (req *Paginate) PaginateScope() func(db *gorm.DB) *gorm.DB {
	return scope.PaginateScope(req.Page, req.PageSize, req.Sort, req.OrderBy)
}

// Response
type Response struct {
	Status int64       `json:"status"`
	Msg    string      `json:"message"`
	Data   interface{} `json:"data"`
}
