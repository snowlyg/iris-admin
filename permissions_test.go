package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPermissions(t *testing.T) {
	url := "/v1/admin/permissions"
	getMore(t, url, iris.StatusOK, true, "操作成功", nil)
}

func TestPermissionCreate(t *testing.T) {

	oj := map[string]interface{}{
		"name":         "test_create_role",
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

func TestPermissionUpdate(t *testing.T) {

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
	update(t, fmt.Sprintf(url, testPerm.ID), oj, iris.StatusOK, true, "操作成功", data)
}

func TestPermissionDelete(t *testing.T) {

	url := "/v1/admin/permissions/%d"
	delete(t, fmt.Sprintf(url, testPerm.ID), iris.StatusOK, true, "删除成功", nil)
}
