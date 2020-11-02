// +build test perm api

package tests

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPermissions(t *testing.T) {
	getMore(t, "permissions", iris.StatusOK, 200, "操作成功")
}

func TestPermissionCreate(t *testing.T) {
	oj := map[string]interface{}{
		"name":         "create_role",
		"display_name": "create_display_name",
		"description":  "create_description",
	}

	create(t, "permissions", oj, iris.StatusOK, 200, "操作成功")
}

func TestPermissionUpdate(t *testing.T) {
	tr, err := CreatePermission("tname1", "tdsiName1", "tdec1")
	if err != nil {
		fmt.Print(err)
	}
	oj := map[string]interface{}{
		"name":         "test_update_role",
		"display_name": "update_display_name",
		"description":  "update_description",
	}

	url := "permissions/%d"
	update(t, fmt.Sprintf(url, tr.ID), oj, iris.StatusOK, 200, "操作成功")
}

func TestPermissionDelete(t *testing.T) {
	tr, err := CreatePermission("tname2", "tdsiName2", "tdec2")
	if err != nil {
		fmt.Print(err)
	}
	delete(t, fmt.Sprintf("permissions/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
