package main

import (
	"IrisApiProject/models"
	"fmt"
	"testing"

	"github.com/kataras/iris"
)

// 后台权限列表
func TestPermissions(t *testing.T) {

	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("permissions")

	// 发起 http 请求
	// Url        string      //测试路由
	// Object     interface{} //发送的json 对象
	// StatusCode int         //返回的 http 状态码
	// Status     bool        //返回的状态
	// Msg        string      //返回提示信息
	// Data       interface{} //返回数据
	url := "/v1/admin/permissions"
	getMore(t, url, iris.StatusOK, true, "操作成功", nil)
}

// 创建权限
func TestPermissionCreate(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("permissions")

	oj := map[string]interface{}{
		"name":         "create_permission",
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

	url := "/v1/admin/permissions"
	create(t, url, oj, iris.StatusOK, true, "操作成功", data)
}

// 更新权限
func TestPermissionUpdate(t *testing.T) {

	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("permissions")

	aul := &models.PermissionJson{
		Name:        "create_permission",
		Description: "create_description",
		DisplayName: "create_display_name",
		Level:       888,
	}

	testAdminPermission := models.CreatePermission(aul)

	oj := map[string]interface{}{
		"name":         "update_permission",
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

	url := "/v1/admin/permissions/%d/update"
	update(t, fmt.Sprintf(url, testAdminPermission.ID), oj, iris.StatusOK, true, "操作成功", data)
}

// 删除权限
func TestPermissionDelete(t *testing.T) {
	// 设置测试数据表
	// 测试前后会自动创建和删除表
	SetTestTableName("permissions")

	aul := new(models.PermissionJson)

	aul.Name = "guest"
	aul.DisplayName = "访客"
	aul.Description = "访客"
	aul.Level = 999

	delPermission := models.CreatePermission(aul)

	url := "/v1/admin/permissions/%d"
	delete(t, fmt.Sprintf(url, delPermission.ID), iris.StatusOK, true, "删除成功", nil)
}
