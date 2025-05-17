package admin

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snowlyg/helper/str"
	"gorm.io/gorm"
)

type ErrMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

var (
	ErrParamValidate      = errors.New("param unvalidate")
	ErrPaginateParam      = errors.New("paginate param unvalidate")
	ErrUnSupportFramework = errors.New("unsupport framework")
)

// BaseModel
type BaseModel struct {
	Id        uint   `json:"id"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
	DeletedAt string `json:"deletedAt"`
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

// SoftDeleteScope
func SoftDeleteScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL")
	}
}

// IdScope
func IdScope(id any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// InIdsScope
func InIdsScope(ids []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", ids)
	}
}

// InNamesScope
func InNamesScope(names []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name in ?", names)
	}
}

// InUuidsScope
func InUuidsScope(uuids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("uuid in ?", uuids)
	}
}

// NeIdScope
func NeIdScope(id any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id != ?", id)
	}
}

// PaginateScope 	return paginate scope for gorm
func PaginateScope(page, pageSize int, sort, orderBy string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize := getPageSize(pageSize)
		offset := getOffset(page, pageSize)
		return db.Order(getOrderBy(sort, orderBy)).Offset(offset).Limit(pageSize)
	}
}

// getOffset
func getOffset(page, pageSize int) int {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	if page < 0 {
		offset = -1
	}
	return offset
}

// getPageSize
func getPageSize(pageSize int) int {
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize < 0:
		pageSize = -1
	case pageSize == 0:
		pageSize = 10
	}
	return pageSize
}

// getOrderBy
func getOrderBy(sort, orderBy string) string {
	if sort == "" {
		sort = "desc"
	}
	if orderBy == "" {
		orderBy = "created_at"
	}
	return str.Join(orderBy, " ", sort)
}

const (
	ResponseOkMessage    = "OK"
	ResponseErrorMessage = "FAIL"
)

type Response struct {
	Code int    `json:"status"`
	Data any    `json:"data,omitempty"`
	Msg  string `json:"message"`
}

func Ok(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  ResponseOkMessage,
	})
}

func OkWithMessage(message string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  message,
	})
}

func OkWithData(data any, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  ResponseOkMessage,
		Data: data,
	})
}

func Fail(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code: http.StatusBadRequest,
		Msg:  ResponseErrorMessage,
	})
}

// func UnauthorizedFailWithMessage(message string, ctx *gin.Context) {
// }

// func UnauthorizedFailWithDetailed(data any, message string, ctx *gin.Context) {
// }

// func ForbiddenFailWithMessage(message string, ctx *gin.Context) {
// }

func FailWithMessage(message string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code: http.StatusBadRequest,
		Msg:  message,
	})
}

func FailWithDetailed(data any, message string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code: http.StatusBadRequest,
		Msg:  message,
		Data: data,
	})
}

type PageResult struct {
	List     any   `json:"list,omitempty"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

type BaseResponse struct {
	Id        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page" validate:"required"`
	PageSize int    `json:"pageSize" form:"pageSize" validate:"required"`
	OrderBy  string `json:"orderBy" form:"orderBy"`
	SortBy   string `json:"sortBy" form:"sortBy"`
}

type IdsBinding struct {
	Ids []uint `json:"ids" form:"ids" validate:"required,dive,required"`
}

// Request get id data form the context of every query
func (req *IdsBinding) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBind(req); err != nil {
		return ErrParamValidate
	}
	return nil
}

// Id the struct has used to get id form the context of every query
type Id struct {
	Id uint `json:"id" uri:"id"`
}

// Request get id data form the context of every query
func (req *Id) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindUri(req); err != nil {
		return ErrParamValidate
	}
	return nil
}

type Empty struct{}
