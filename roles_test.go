package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
)

// 后台账号列表
func TestRoles(t *testing.T) {

	url := "/v1/admin/roles"
	getMore(t, url, iris.StatusOK, true, "操作成功", nil)
}

// 创建角色
func TestRoleCreate(t *testing.T) {

	oj := map[string]interface{}{
		"name":         "create_role",
		"display_name": "create_display_name",
		"description":  "create_description",
	}

	url := "/v1/admin/roles"
	create(t, url, oj, iris.StatusOK, true, "操作成功", nil)
}

// 更新角色
func TestRoleUpdate(t *testing.T) {
	tr := CreateRole("tname1", "tdsiName1", "tdec1")
	oj := map[string]interface{}{
		"name":         "test_update_role",
		"display_name": "update_display_name",
		"description":  "update_description",
	}

	url := "/v1/admin/roles/%d"
	update(t, fmt.Sprintf(url, tr.ID), oj, iris.StatusOK, true, "操作成功", nil)
}

// 删除角色
func TestRoleDelete(t *testing.T) {
	tr := CreateRole("tname2", "tdsiName2", "tdec2")
	url := "/v1/admin/roles/%d"
	delete(t, fmt.Sprintf(url, tr.ID), iris.StatusOK, true, "删除成功", nil)
}
