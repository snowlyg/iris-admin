package main

import (
	"IrisApiProject/config"
	"IrisApiProject/models"
	"fmt"
	"testing"

	"github.com/kataras/iris"
)

// 后台账号列表
func TestUsers(t *testing.T) {

	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	// 发起 http 请求
	// Url        string      //测试路由
	// Object     interface{} //发送的json 对象
	// StatusCode int         //返回的 http 状态码
	// Status     bool        //返回的状态
	// Msg        string      //返回提示信息
	// Data       interface{} //返回数据
	getMore(t, "/v1/admin/users", iris.StatusOK, true, "操作成功", nil)
}

// 登陆用户信息
func TestUserProfile(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	getMore(t, "/v1/admin/users/profile", iris.StatusOK, true, "操作成功", nil)
}

// 创建用户
func TestUserCreate(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": "test_user",
		"password": "password",
		"name":     "name",
	}

	data := map[string]interface{}{
		"Username": "test_user",
		"Name":     "name",
	}

	create(t, "/v1/admin/users", oj, iris.StatusOK, true, "操作成功", data)
}

// 更新用户
func TestUserUpdate(t *testing.T) {

	testAdminUser := models.CreaterSystemAdmin()
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": "update_user",
		"password": "update_name",
		"name":     "update_name",
	}

	data := map[string]interface{}{
		"Username": oj["username"],
		"Name":     oj["name"],
	}

	update(t, fmt.Sprintf("/v1/admin/users/%d/update", testAdminUser.ID), oj, iris.StatusOK, true, "操作成功", data)
}

// 删除用户
func TestUserDelete(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	aul := new(models.UserJson)
	aul.Username = config.Conf.Get("test.LoginUserName").(string)
	aul.Password = config.Conf.Get("test.LoginPwd").(string)
	aul.Name = config.Conf.Get("test.LoginName").(string)

	delUser := models.CreateUser(aul)

	delete(t, fmt.Sprintf("/v1/admin/users/%d", delUser.ID), iris.StatusOK, true, "删除成功", nil)
}
