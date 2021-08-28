package role

import "github.com/snowlyg/iris-admin/g"

type Request struct {
	BaseRole
}

type ReqPaginate struct {
	g.Paginate
	Name string `json:"name"`
}
