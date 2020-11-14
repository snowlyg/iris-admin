// +build test role api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/tests/mock"
	"testing"
)

// TestRoles
func TestRoles(t *testing.T) {
	tr, err := CreateRole()
	if err != nil {
		color.Red("TestRoles %+v", err)
		return
	}

	obj := map[string]interface{}{"limit": 1, "page": 1}
	more := &More{tr.ID, 1, 1, 2}
	getMore(t, "roles", iris.StatusOK, obj, more)
}

func TestRolesNoPagination(t *testing.T) {
	tr, err := CreateRole()
	if err != nil {
		color.Red("TestRoles %+v", err)
		return
	}

	obj := map[string]interface{}{"limit": -1, "page": -1}
	more := &More{tr.ID, -1, 3, 3}
	getMore(t, "roles", iris.StatusOK, obj, more)
}

// TestRoleCreate
func TestRoleCreate(t *testing.T) {
	m := mock.Role{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestRoleCreate %+v", err)
		return
	}
	create(t, "roles", m, iris.StatusOK, 200, "操作成功")
}

// TestRoleUpdate
func TestRoleUpdate(t *testing.T) {
	m := mock.Role{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestRoleUpdate %+v", err)
		return
	}

	tr, err := CreateRole()
	if err != nil {
		color.Red("TestRoleUpdate %+v", err)
		return
	}

	url := "roles/%d"
	update(t, fmt.Sprintf(url, tr.ID), m, iris.StatusOK, 200, "操作成功")
}

// 删除角色
func TestRoleDelete(t *testing.T) {
	tr, err := CreateRole()
	if err != nil {
		fmt.Print(err)
	}
	delete(t, fmt.Sprintf("roles/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
