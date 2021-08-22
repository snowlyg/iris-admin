package perm

import "github.com/snowlyg/iris-admin/g"

type Request struct {
	BasePerission
}

type ReqPaginate struct {
	g.Paginate
	Name string `json:"name"`
}
