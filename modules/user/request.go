package user

import "github.com/snowlyg/iris-admin/g"

type Request struct {
	BaseUser
}

type ReqPaginate struct {
	g.Paginate
	Name string `json:"name"`
}
