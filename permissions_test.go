package main

import (
	"IrisApiProject/models"
	"fmt"
	"testing"

	"github.com/kataras/iris"
)

// 后台权限列表
func TestPermissions(t *testing.T) {

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

	oj := map[string]interface{}{
		"name":         "test_create_permission",
		"display_name": "create_display_name",
		"description":  "create_description",
	}

	data := map[string]interface{}{
		"Name":        oj["name"],
		"DisplayName": oj["display_name"],
		"Description": oj["description"],
	}

	url := "/v1/admin/permissions"
	create(t, url, oj, iris.StatusOK, true, "操作成功", data)
}

// 更新权限
func TestPermissionUpdate(t *testing.T) {

	aul := &models.PermissionJson{
		Name:        "test_update_permission",
		Description: "create_description",
		DisplayName: "create_display_name",
	}

	testAdminPermission := models.CreatePermission(aul)

	oj := map[string]interface{}{
		"name":         "update_permission",
		"display_name": "update_display_name",
		"description":  "update_description",
	}

	data := map[string]interface{}{
		"Name":        oj["name"],
		"DisplayName": oj["display_name"],
		"Description": oj["description"],
	}

	url := "/v1/admin/permissions/%d"
	update(t, fmt.Sprintf(url, testAdminPermission.ID), oj, iris.StatusOK, true, "操作成功", data)
}

// 删除权限
func TestPermissionDelete(t *testing.T) {

	aul := new(models.PermissionJson)

	aul.Name = "guest"
	aul.DisplayName = "访客"
	aul.Description = "访客"

	delPermission := models.CreatePermission(aul)

	url := "/v1/admin/permissions/%d"
	delete(t, fmt.Sprintf(url, delPermission.ID), iris.StatusOK, true, "删除成功", nil)
}
