package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ErrMsg struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

var (
	ErrParamValidate      = errors.New("param unvalidate")
	ErrPaginateParam      = errors.New("paginate param unvalidate")
	ErrUnSupportFramework = errors.New("unsupport framework")
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
func (req *ReqId) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindUri(req); err != nil {
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

func (req *Paginate) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBind(req); err != nil {
		return ErrParamValidate
	}
	return nil
}

// PaginateScope paginate scope
func (req *Paginate) PaginateScope() func(db *gorm.DB) *gorm.DB {
	return PaginateScope(req.Page, req.PageSize, req.Sort, req.OrderBy)
}
