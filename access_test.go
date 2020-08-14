// +build test

package IrisAdminApi

import (
	"testing"

	"github.com/snowlyg/IrisAdminApi/config"
)

// 登陆成功
func TestUserLoginSuccess(t *testing.T) {
	oj := map[string]string{
		"username": config.Config.Admin.UserName,
		"password": config.Config.Admin.Pwd,
	}
	login(t, oj, iris.StatusOK, 200, "登陆成功")
}

// 输入不存在的用户名登陆
func TestUserLoginWithErrorName(t *testing.T) {
	oj := map[string]string{
		"username": "err_user",
		"password": config.Config.Admin.Pwd,
	}

	login(t, oj, iris.StatusOK, 400, "用户不存在")
}

// 输入错误的登陆密码
func TestUserLoginWithErrorPwd(t *testing.T) {

	oj := map[string]string{
		"username": config.Config.Admin.UserName,
		"password": "admin",
	}
	login(t, oj, iris.StatusOK, 400, "用户名或密码错误")
}

// 不输入用户名
func TestUserLoginWithNoUsername(t *testing.T) {

	oj := map[string]string{
		"username": "",
		"password": "admin",
	}
	login(t, oj, iris.StatusOK, 400, "用户名为必填字段")
}

// 不输入密码
func TestUserLoginWithNoPwd(t *testing.T) {

	oj := map[string]string{
		"username": "username",
		"password": "",
	}
	login(t, oj, iris.StatusOK, 400, "密码为必填字段")
}

// 输入登陆密码格式错误
func TestUserLoginWithErrorFormtPwd(t *testing.T) {
	oj := map[string]string{
		"username": config.Config.Admin.UserName,
		"password": "123",
	}

	login(t, oj, iris.StatusOK, 400, "用户名或密码错误")
}

// 输入登陆密码格式错误
func TestUserLoginWithErrorFormtUserName(t *testing.T) {

	oj := map[string]string{
		"username": "df",
		"password": "123",
	}

	login(t, oj, iris.StatusOK, 400, "用户不存在")
}
