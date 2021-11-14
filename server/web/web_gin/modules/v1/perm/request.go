package perm

import (
	"github.com/snowlyg/iris-admin/server/database/orm"
)

type DeleteApiReq struct {
	Id     uint   `json:"id" form:"id" binding:"required,gt=0"`
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
}

type AuthorityType struct {
	AuthorityType int `json:"authorityType" form:"authorityType"`
}

type ReqPaginate struct {
	orm.Paginate
	Path        string `json:"path" form:"path"`
	Description string `json:"description" form:"description"`
	ApiGroup    string `json:"apiGroup" form:"apiGroup"`
	Method      string `json:"method" form:"method"`
	OrderKey    string `json:"orderKey" form:"orderKey"`
	Desc        bool   `json:"desc" form:"desc"`
}
