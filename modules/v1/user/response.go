package user

import (
	"regexp"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
)

type Response struct {
	g.Model
	BaseUser
	Roles []string `gorm:"-" json:"roles"`
}

func (res *Response) ToString() {
	if res.Avatar == "" {
		return
	}
	re := regexp.MustCompile("^http")
	if !re.MatchString(res.Avatar) {
		res.Avatar = str.Join("http://127.0.0.1:8085/upload/", res.Avatar)
	}
}

type LoginResponse struct {
	g.ReqId
	Password string `json:"password"`
}
