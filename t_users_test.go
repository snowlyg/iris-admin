package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris"
)

// 后台账号列表
func TestUsersList(t *testing.T) {

	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	// users := []*Users{testAdminUser}

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
		"phone":    "13412334567",
		"role_id":  1,
	}

	data := map[string]interface{}{
		"Username": "test_user",
		"Name":     "name",
		"Phone":    "13412334567",
		"RoleId":   1,
	}

	create(t, "/v1/admin/users", oj, iris.StatusOK, true, "操作成功", data)
}

// 创建用户
func TestUserUpdate(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": "update_user",
		"password": "update_name",
		"name":     "update_name",
		"phone":    "13412334567",
		"role_id":  2,
	}

	data := map[string]interface{}{
		"Username": oj["username"],
		"Name":     oj["name"],
		"Phone":    oj["phone"],
		"RoleId":   oj["role_id"],
	}

	update(t, fmt.Sprintf("/v1/admin/users/%d/update", testAdminUser.ID), oj, iris.StatusOK, true, "操作成功", data)
}

// 创建用户
func TestUserDelete(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("users")

	oj := map[string]interface{}{
		"username": "update_user",
		"password": "update_name",
		"name":     "update_name",
		"phone":    "13412334567",
		"role_id":  2,
	}

	delete(t, fmt.Sprintf("/v1/admin/users/%d/delete", testAdminUser.ID), iris.StatusOK, true, "操作成功", nil)
}
