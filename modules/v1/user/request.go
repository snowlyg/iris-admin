package user

import "github.com/snowlyg/iris-admin/g"

type Request struct {
	BaseUser
	Password string `json:"password"`
	RoleIds  []uint `json:"role_ids"`
}

type ReqPaginate struct {
	g.Paginate
	Name string `json:"name"`
}
