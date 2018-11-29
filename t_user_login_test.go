package main

import (
	"github.com/kataras/iris"
	"testing"
)

//登陆成功
func TestAdminLoginSuccess(t *testing.T) {
	//设置测试数据表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": conf.Get("test.LoginUserName").(string),
		"password": conf.Get("test.LoginPwd").(string),
	}
	bc := BaseCase{"/v1/admin/login", oj, iris.StatusOK, true, "登陆成功", nil}
	bc.login(t)
}

// 输入不存在的用户名登陆
func TestUserLoginWithErrorName(t *testing.T) {
	//设置测试数据表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": "err_user",
		"password": conf.Get("test.LoginPwd").(string),
	}
	bc := BaseCase{"/v1/admin/login", oj, iris.StatusOK, false, "用户不存在", nil}
	bc.login(t)
}

//输入错误的登陆密码
func TestUserLoginWithErrorPwd(t *testing.T) {
	//设置测试数据表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": conf.Get("test.LoginUserName").(string),
		"password": "admin",
	}
	bc := BaseCase{"/v1/admin/login", oj, iris.StatusOK, false, "用户名或密码错误", nil}
	bc.login(t)
}

//输入登陆密码格式错误
func TestUserLoginWithErrorFormtPwd(t *testing.T) {
	//设置测试数据表
	//测试前后会自动创建和删除表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": conf.Get("test.LoginUserName").(string),
		"password": "123",
	}
	bc := BaseCase{"/v1/admin/login", oj, iris.StatusOK, false, "密码格式错误", nil}

	bc.login(t)
}

//输入登陆密码格式错误
func TestUserLoginWithErrorFormtUserName(t *testing.T) {
	//设置测试数据表
	//测试前后会自动创建和删除表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": "df",
		"password": "123",
	}
	bc := BaseCase{"/v1/admin/login", oj, iris.StatusOK, false, "用户名格式错误", nil}
	bc.login(t)
}
