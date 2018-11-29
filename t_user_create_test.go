package main

import (
	"github.com/kataras/iris"
	"testing"
)

//登陆用户信息
func TestUserCreate(t *testing.T) {
	//设置测试数据表
	//测试前后会自动创建和删除表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": conf.Get("test.LoginUserName").(string),
		"password": conf.Get("test.LoginPwd").(string),
		"name":     "name",
		"phone":    "13412334567",
		"RoleId":   1,
	}

	bc := BaseCase{"/v1/admin/users", oj, iris.StatusOK, true, "操作成功", oj}
	bc.post(t)
}
