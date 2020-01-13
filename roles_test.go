package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
)

// 后台账号列表
func TestRoles(t *testing.T) {
	getMore(t, baseUrl+"roles", iris.StatusOK, true, "操作成功")
}

// 创建角色
func TestRoleCreate(t *testing.T) {

	oj := map[string]interface{}{
		"name":         "create_role",
		"display_name": "create_display_name",
		"description":  "create_description",
	}

	create(t, baseUrl+"roles", oj, iris.StatusOK, true, "操作成功")
}

// 更新角色
func TestRoleUpdate(t *testing.T) {
	tr := CreateRole("tname1", "tdsiName1", "tdec1")
	oj := map[string]interface{}{
		"name":         "test_update_role",
		"display_name": "update_display_name",
		"description":  "update_description",
	}

	url := baseUrl + "roles/%d"
	update(t, fmt.Sprintf(url, tr.ID), oj, iris.StatusOK, true, "操作成功")
}

// 删除角色
func TestRoleDelete(t *testing.T) {
	tr := CreateRole("tname2", "tdsiName2", "tdec2")
	delete(t, fmt.Sprintf(baseUrl+"roles/%d", tr.ID), iris.StatusOK, true, "删除成功")
}
