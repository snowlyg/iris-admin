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
	Id        uint   `json:"id"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
	DeletedAt string `json:"deletedAt"`
}

// UriId the struct has used to get id form the context of every query
type UriId struct {
	Id uint `json:"id" uri:"id"`
}

// Request get id data form the context of every query
func (req *UriId) Request(ctx *gin.Context) error {
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

// IdScope
// - id uint
func IdScope(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// InIdsScope
// - ids []uint
func InIdsScope(ids []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", ids)
	}
}

// InNamesScope
// - names []string
func InNamesScope(names []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name in ?", names)
	}
}

// InUuidsScope
// - uuids []string
func InUuidsScope(uuids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("uuid in ?", uuids)
	}
}

// NeIdScope
// - id uint
func NeIdScope(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id != ?", id)
	}
}

// PaginateScope 	return paginate scope for gorm
// - page 			int
// - pageSize 	int
// - sort 			string
// - orderBy 		string
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
	Code int         `json:"status"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"message"`
}

func Result(code int, data interface{}, msg string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{code, data, msg})
}

func Ok(ctx *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, ResponseOkMessage, ctx)
}

func OkWithMessage(message string, ctx *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, message, ctx)
}

func OkWithData(data interface{}, ctx *gin.Context) {
	Result(http.StatusOK, data, ResponseOkMessage, ctx)
}

func OkWithDetailed(data interface{}, message string, ctx *gin.Context) {
	Result(http.StatusOK, data, message, ctx)
}

func Fail(ctx *gin.Context) {
	Result(http.StatusBadRequest, map[string]interface{}{}, ResponseErrorMessage, ctx)
}

func UnauthorizedFailWithMessage(message string, ctx *gin.Context) {
	Result(http.StatusUnauthorized, map[string]interface{}{}, message, ctx)
}

func UnauthorizedFailWithDetailed(data interface{}, message string, ctx *gin.Context) {
	Result(http.StatusUnauthorized, data, message, ctx)
}

func ForbiddenFailWithMessage(message string, ctx *gin.Context) {
	Result(http.StatusForbidden, map[string]interface{}{}, message, ctx)
}

func FailWithMessage(message string, ctx *gin.Context) {
	Result(http.StatusBadRequest, map[string]interface{}{}, message, ctx)
}

func FailWithDetailed(data interface{}, message string, ctx *gin.Context) {
	Result(http.StatusBadRequest, data, message, ctx)
}

type PageResult struct {
	List     interface{} `json:"list,omitempty"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type BaseResponse struct {
	Id        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page" binding:"required"`
	PageSize int    `json:"pageSize" form:"pageSize" binding:"required"`
	OrderBy  string `json:"orderBy" form:"orderBy"`
	SortBy   string `json:"sortBy" form:"sortBy"`
}

type IdsBinding struct {
	Ids []uint `json:"ids" form:"ids" binding:"required,dive,required"`
}

type Empty struct{}
