package role

import "github.com/snowlyg/iris-admin/g"

type Request struct {
	BaseRole
	Perms [][]string `json:"perms"`
}

type ReqPaginate struct {
	g.Paginate
	Name string `json:"name"`
}
