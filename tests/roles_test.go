// +build test role api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/tests/mock"
	"testing"
)

// 后台账号列表
func TestRoles(t *testing.T) {
	getMore(t, "roles", iris.StatusOK, 200, "操作成功")
}

// 创建角色
func TestRoleCreate(t *testing.T) {
	m := mock.Role{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestRoleCreate %+v", err)
		return
	}
	create(t, "roles", m, iris.StatusOK, 200, "操作成功")
}

// 更新角色
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
