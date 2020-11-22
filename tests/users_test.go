// +build test user api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/tests/mock"
	"testing"
)

func TestUsers(t *testing.T) {
	tu, err := CreateUser()
	if err != nil {
		color.Red("TestUserUpdate %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": 1, "page": 1, "field": "id,name,username,created_at"}
	more := &More{tu.ID, 1, 1, 2, []interface{}{"id", "name", "username", "roles", "created_at"}}
	getMore(t, "users", iris.StatusOK, obj, more)
}

func TestUsersNoPagination(t *testing.T) {
	tu, err := CreateUser()
	if err != nil {
		color.Red("TestUsersNoPagination %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": -1, "page": -1, "field": "id,name,username,created_at"}
	more := &More{tu.ID, -1, 3, 3, []interface{}{"id", "name", "username", "roles", "created_at"}}
	getMore(t, "users", iris.StatusOK, obj, more)
}

func TestUserProfile(t *testing.T) {
	obj := map[string]interface{}{"limit": 1, "page": 1}
	getData(t, "profile", iris.StatusOK, obj, []interface{}{"avatar", "id", "created_at", "introduction", "roles", "role_ids", "name"})
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
