package main

import (
	"IrisApiProject/models"
	"fmt"
	"testing"

	"github.com/kataras/iris"
)

// 后台账号列表
func TestRoles(t *testing.T) {

	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("roles")

	// 发起 http 请求
	// Url        string      //测试路由
	// Object     interface{} //发送的json 对象
	// StatusCode int         //返回的 http 状态码
	// Status     bool        //返回的状态
	// Msg        string      //返回提示信息
	// Data       interface{} //返回数据
	getMore(t, "/v1/admin/roles", iris.StatusOK, true, "操作成功", nil)
}

// 创建角色
func TestRoleCreate(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("roles")

	oj := map[string]interface{}{
		"name":         "create_role",
		"display_name": "create_display_name",
		"description":  "create_description",
		"level":        888,
	}

	data := map[string]interface{}{
		"Name":        oj["name"],
		"DisplayName": oj["display_name"],
		"Description": oj["description"],
		"Level":       oj["level"],
	}

	create(t, "/v1/admin/roles", oj, iris.StatusOK, true, "操作成功", data)
}

// 更新角色
func TestRoleUpdate(t *testing.T) {

	testAdminRole := models.CreaterSystemAdminRole()
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("roles")

	oj := map[string]interface{}{
		"name":         "update_role",
		"display_name": "update_display_name",
		"description":  "update_description",
		"level":        888,
	}

	data := map[string]interface{}{
		"Name":        oj["name"],
		"DisplayName": oj["display_name"],
		"Description": oj["description"],
		"Level":       oj["level"],
	}

	update(t, fmt.Sprintf("/v1/admin/roles/%d/update", testAdminRole.ID), oj, iris.StatusOK, true, "操作成功", data)
}

// 删除角色
func TestRoleDelete(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("roles")

	aul := new(models.RoleJson)

	aul.Name = "guest"
	aul.DisplayName = "访客"
	aul.Description = "访客"
	aul.Level = 999

	delRole := models.CreateRole(aul)

	delete(t, fmt.Sprintf("/v1/admin/roles/%d", delRole.ID), iris.StatusOK, true, "删除成功", nil)
}
