package request

// Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page" validate:"required"`
	PageSize int    `json:"pageSize" form:"pageSize" validate:"required"`
	OrderBy  string `json:"orderBy" form:"orderBy"`
	SortBy   string `json:"sortBy" form:"sortBy"`
}

// Find by id structure
type IdBinding struct {
	Id uint `json:"id" uri:"id" validate:"required"`
}

type IdsBinding struct {
	Ids []uint `json:"ids" form:"ids" validate:"required,dive,required"`
}

type Empty struct{}
