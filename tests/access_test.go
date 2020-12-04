// +build test access api

package tests

import (
	"github.com/kataras/iris/v12"
	"testing"

	"github.com/snowlyg/blog/libs"
)

func TestUserLoginSuccess(t *testing.T) {
	oj := map[string]string{
		"username": libs.Config.Admin.UserName,
		"password": libs.Config.Admin.Pwd,
	}
	login(t, oj, iris.StatusOK, 200, "请求成功")
}

func TestUserLoginWithErrorName(t *testing.T) {
	oj := map[string]string{
		"username": "err_user",
		"password": libs.Config.Admin.Pwd,
	}

	login(t, oj, iris.StatusOK, 5004, "数据为空")
}

func TestUserLoginWithErrorPwd(t *testing.T) {

	oj := map[string]string{
		"username": libs.Config.Admin.UserName,
		"password": "admin",
	}
	login(t, oj, iris.StatusOK, 400, "用户名或密码错误")
}

func TestUserLoginWithNoUsername(t *testing.T) {

	oj := map[string]string{
		"username": "",
		"password": "admin",
	}
	login(t, oj, iris.StatusOK, 400, "用户名为必填字段")
}

func TestUserLoginWithNoPwd(t *testing.T) {
	oj := map[string]string{
		"username": "username",
		"password": "",
	}
	login(t, oj, iris.StatusOK, 400, "密码为必填字段")
}

func TestUserLoginWithErrorFormtPwd(t *testing.T) {
	oj := map[string]string{
		"username": libs.Config.Admin.UserName,
		"password": "123",
	}

	login(t, oj, iris.StatusOK, 400, "用户名或密码错误")
}

func TestUserLoginWithErrorFormtUserName(t *testing.T) {
	oj := map[string]string{
		"username": "df",
		"password": "123",
	}
	login(t, oj, iris.StatusOK, 5004, "数据为空")
}
