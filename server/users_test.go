// +build test

package main

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/kataras/iris/v12"
)

func TestUsers(t *testing.T) {
	getMore(t, "users", iris.StatusOK, 200, "操作成功")
}

func TestUserProfile(t *testing.T) {
	getMore(t, "profile", iris.StatusOK, 200, "")
}

func TestUserCreate(t *testing.T) {
	tr := CreateRole("tname3", "tdsiName", "tdec")
	oj := map[string]interface{}{
		"username": gofakeit.Name(),
		"password": "password",
		"name":     "name",
		"role_ids": []uint{tr.ID},
	}

	create(t, "users", oj, iris.StatusOK, 200, "操作成功")
}

func TestUserUpdate(t *testing.T) {
	tr := CreateRole("tname4", "tdsiName", "tdec")
	oj := map[string]interface{}{
		"username": gofakeit.Name(),
		"password": "update_name",
		"name":     "update_name",
		"role_ids": []uint{tr.ID},
	}

	tu := CreateUser()
	update(t, fmt.Sprintf("users/%d", tu.ID), oj, iris.StatusOK, 200, "操作成功")
}

func TestUserDelete(t *testing.T) {
	tu := CreateUser()
	delete(t, fmt.Sprintf("users/%d", tu.ID), iris.StatusOK, 200, "删除成功")
}
