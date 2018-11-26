package main

import (
	"github.com/kataras/iris"
	"testing"
)

//后台账号列表
func TestUsersList(t *testing.T) {
	bc := BaseCase{"/v1/admin/users", nil, iris.StatusOK, true, "操作成功", nil}
	bc.get(t)
}
