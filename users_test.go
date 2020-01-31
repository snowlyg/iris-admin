package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestUsers(t *testing.T) {
	getMore(t, baseUrl+"users", iris.StatusOK, true, "操作成功")
}

func TestUserProfile(t *testing.T) {
	getMore(t, baseUrl+"users/profile", iris.StatusOK, true, "")
}

func TestUserCreate(t *testing.T) {
	tr := CreateRole("tname3", "tdsiName", "tdec")
	oj := map[string]interface{}{
		"username": "test_user",
		"password": "password",
		"name":     "name",
		"role_ids": []uint{tr.ID},
	}

	create(t, baseUrl+"users", oj, iris.StatusOK, true, "操作成功")
}

func TestUserUpdate(t *testing.T) {
	tr := CreateRole("tname4", "tdsiName", "tdec")
	oj := map[string]interface{}{
		"username": "test_update_user",
		"password": "update_name",
		"name":     "update_name",
		"role_ids": []uint{tr.ID},
	}

	tu := CreateUser()
	update(t, fmt.Sprintf(baseUrl+"users/%d", tu.ID), oj, iris.StatusOK, true, "操作成功")
}

func TestUserDelete(t *testing.T) {
	tu := CreateUser()
	delete(t, fmt.Sprintf(baseUrl+"users/%d", tu.ID), iris.StatusOK, true, "删除成功")
}
