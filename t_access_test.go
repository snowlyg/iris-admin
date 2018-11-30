package main

import (
	"testing"

	"github.com/kataras/iris"
)

// 登陆成功
func TestAdminLoginSuccess(t *testing.T) {
	// 设置测试数据表
	SetTestTableName("users")

	oj := map[string]string{
		"username": conf.Get("test.LoginUserName").(string),
		"password": conf.Get("test.LoginPwd").(string),
	}

	login(t, "/v1/admin/login", oj, iris.StatusOK, true, "登陆成功", nil)
}

// 输入不存在的用户名登陆
func TestUserLoginWithErrorName(t *testing.T) {
	// 设置测试数据表
	SetTestTableName("users")

	oj := map[string]string{
		"username": "err_user",
		"password": conf.Get("test.LoginPwd").(string),
	}

	login(t, "/v1/admin/login", oj, iris.StatusOK, false, "用户不存在", nil)
}

// 输入错误的登陆密码
func TestUserLoginWithErrorPwd(t *testing.T) {
	// 设置测试数据表
	SetTestTableName("users")

	oj := map[string]string{
		"username": conf.Get("test.LoginUserName").(string),
		"password": "admin",
	}
	login(t, "/v1/admin/login", oj, iris.StatusOK, false, "用户名或密码错误", nil)
}

// 输入登陆密码格式错误
func TestUserLoginWithErrorFormtPwd(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	oj := map[string]string{
		"username": conf.Get("test.LoginUserName").(string),
		"password": "123",
	}

	login(t, "/v1/admin/login", oj, iris.StatusOK, false, "密码格式错误", nil)
}

// 输入登陆密码格式错误
func TestUserLoginWithErrorFormtUserName(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	oj := map[string]string{
		"username": "df",
		"password": "123",
	}

	login(t, "/v1/admin/login", oj, iris.StatusOK, false, "用户名格式错误", nil)
}
