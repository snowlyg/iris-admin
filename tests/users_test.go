// +build test user api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/tests/mock"
	"testing"
)

func TestUsers(t *testing.T) {
	getMore(t, "users", iris.StatusOK, 200, "操作成功")
}

func TestUserProfile(t *testing.T) {
	getMore(t, "profile", iris.StatusOK, 200, "请求成功")
}

func TestUserCreate(t *testing.T) {
	m := mock.User{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestUserCreate %+v", err)
		return
	}

	create(t, "users", m, iris.StatusOK, 200, "操作成功")
}

func TestUserUpdate(t *testing.T) {

	m := mock.User{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestUserUpdate %+v", err)
		return
	}

	tu, err := CreateUser()
	if err != nil {
		color.Red("TestUserUpdate %+v", err)
		return
	}
	update(t, fmt.Sprintf("users/%d", tu.ID), m, iris.StatusOK, 200, "操作成功")
}

func TestUserDelete(t *testing.T) {
	tu, err := CreateUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	delete(t, fmt.Sprintf("users/%d", tu.ID), iris.StatusOK, 200, "删除成功")
}
