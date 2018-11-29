package main

import (
	"github.com/kataras/iris"
	"testing"
)

//登陆用户信息
func TestUserProfile(t *testing.T) {
	//设置测试数据表
	//测试前后会自动创建和删除表
	SetTestTableName("users")

	bc := BaseCase{"/v1/admin/users/profile", nil, iris.StatusOK, true, "操作成功", testAdminUser.Name}
	bc.get(t)
}
