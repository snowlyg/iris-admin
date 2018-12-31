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

	url := "/v1/admin/users/profile"
	getMore(t, url, iris.StatusOK, true, "操作成功", nil)
}

// 创建用户
func TestUserCreate(t *testing.T) {

	aul := &models.RoleJson{
		Name:        "test_guest",
		Description: "访客",
		DisplayName: "访客",
	}

	testAdminRole := models.CreateRole(aul)

	oj := map[string]interface{}{
		"username": "test_user",
		"password": "password",
		"name":     "name",
		"role_id":  testAdminRole.ID,
	}

	data := map[string]interface{}{
		"Username": "test_user",
		"Name":     "name",
		"role_id":  testAdminRole.ID,
	}

	url := "/v1/admin/users"
	create(t, url, oj, iris.StatusOK, true, "操作成功", data)
}

// 更新用户
func TestUserUpdate(t *testing.T) {

	role_json := &models.RoleJson{
		Name:        "test_guest",
		Description: "访客",
		DisplayName: "访客",
	}

	testAdminRole := models.CreateRole(role_json)

	oj := map[string]interface{}{
		"username": "test_update_user",
		"password": "update_name",
		"name":     "update_name",
		"role_id":  testAdminRole.ID,
	}

	data := map[string]interface{}{
		"Username": oj["username"],
		"Name":     oj["name"],
	}

	aul := &models.UserJson{
		Username: "guest",
		Name:     "访客",
		Password: "guest111",
		RoleID:   testAdminRole.ID,
	}

	testAdminUser := models.CreateUser(aul)

	url := "/v1/admin/users/%d"
	update(t, fmt.Sprintf(url, testAdminUser.ID), oj, iris.StatusOK, true, "操作成功", data)
}

// 删除用户
func TestUserDelete(t *testing.T) {

	aul := new(models.UserJson)
	aul.Username = config.Conf.Get("test.LoginUserName").(string)
	aul.Password = config.Conf.Get("test.LoginPwd").(string)
	aul.Name = config.Conf.Get("test.LoginName").(string)

	delUser := models.CreateUser(aul)

	url := "/v1/admin/users/%d"
	delete(t, fmt.Sprintf(url, delUser.ID), iris.StatusOK, true, "删除成功", nil)
}
