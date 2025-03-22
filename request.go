package admin

// Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page" binding:"required"`
	PageSize int    `json:"pageSize" form:"pageSize" binding:"required"`
	OrderBy  string `json:"orderBy" form:"orderBy"`
	SortBy   string `json:"sortBy" form:"sortBy"`
}

// Find by id structure
type IdBinding struct {
	Id uint `json:"id" uri:"id" binding:"required"`
}

type IdsBinding struct {
	Ids []uint `json:"ids" form:"ids" binding:"required,dive,required"`
}

type Empty struct{}
