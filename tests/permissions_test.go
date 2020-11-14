// +build test perm api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/tests/mock"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPermissions(t *testing.T) {
	obj := map[string]interface{}{"limit": 1, "page": 1}
	more := &More{53, 1, 1, 53}
	getMore(t, "permissions", iris.StatusOK, obj, more)
}

func TestPermissionsNoPagination(t *testing.T) {
	obj := map[string]interface{}{"limit": -1, "page": -1}
	more := &More{53, -1, 53, 53}
	getMore(t, "permissions", iris.StatusOK, obj, more)
}

func TestPermissionCreate(t *testing.T) {
	m := mock.Permission{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestPermissionCreate %+v", err)
		return
	}

	create(t, "permissions", m, iris.StatusOK, 200, "操作成功")
}

func TestPermissionUpdate(t *testing.T) {
	m := mock.Permission{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestPermissionUpdate %+v", err)
		return
	}

	tr, err := CreatePermission()
	if err != nil {
		fmt.Print(err)
	}

	url := "permissions/%d"
	update(t, fmt.Sprintf(url, tr.ID), m, iris.StatusOK, 200, "操作成功")
}

func TestPermissionDelete(t *testing.T) {
	tr, err := CreatePermission()
	if err != nil {
		fmt.Print(err)
	}
	delete(t, fmt.Sprintf("permissions/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
