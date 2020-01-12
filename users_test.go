package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestUsers(t *testing.T) {
	getMore(t, "/v1/admin/users", iris.StatusOK, true, "操作成功", nil)
}

func TestUserProfile(t *testing.T) {
	url := "/v1/admin/users/profile"
	getMore(t, url, iris.StatusOK, true, "操作成功", nil)
}

func TestUserCreate(t *testing.T) {
	tr := CreateRole("tname3", "tdsiName", "tdec")
	oj := map[string]interface{}{
		"username": "test_user",
		"password": "password",
		"name":     "name",
		"role_ids": []uint{tr.ID},
	}

	url := "/v1/admin/users"
	create(t, url, oj, iris.StatusOK, true, "操作成功", nil)
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
	update(t, fmt.Sprintf("/v1/admin/users/%d", tu.ID), oj, iris.StatusOK, true, "操作成功", nil)
}

func TestUserDelete(t *testing.T) {
	tu := CreateUser()
	delete(t, fmt.Sprintf("/v1/admin/users/%d", tu.ID), iris.StatusOK, true, "删除成功", nil)
}
