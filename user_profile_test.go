package main

import (
	"github.com/kataras/iris"
	"testing"
)

//登陆用户信息
func TestUserProfile(t *testing.T) {
	bc := BaseCase{"/v1/admin/users/profile", nil, iris.StatusOK, true, "操作成功", "Super Admin"}
	bc.get(t)
}
