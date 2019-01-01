package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris"
)

func TestUsers(t *testing.T) {
	getMore(t, "/v1/admin/users", iris.StatusOK, true, "操作成功", nil)
}

func TestUserProfile(t *testing.T) {
	url := "/v1/admin/users/profile"
	getMore(t, url, iris.StatusOK, true, "操作成功", nil)
}

func TestUserCreate(t *testing.T) {

	oj := map[string]interface{}{
		"username": "test_user",
		"password": "password",
		"name":     "name",
		"role_id":  testRole.ID,
	}

	data := map[string]interface{}{
		"Username": oj["username"],
		"Name":     oj["name"],
		"RoleID":   oj["role_id"],
	}

	url := "/v1/admin/users"
	create(t, url, oj, iris.StatusOK, true, "操作成功", data)
}

func TestUserUpdate(t *testing.T) {

	oj := map[string]interface{}{
		"username": "test_update_user",
		"password": "update_name",
		"name":     "update_name",
		"role_id":  testRole.ID,
	}

	data := map[string]interface{}{
		"Username": oj["username"],
		"Name":     oj["name"],
	}

	url := "/v1/admin/users/%d"
	update(t, fmt.Sprintf(url, testUser.ID), oj, iris.StatusOK, true, "操作成功", data)
}

func TestUserDelete(t *testing.T) {

	url := "/v1/admin/users/%d"
	delete(t, fmt.Sprintf(url, testUser.ID), iris.StatusOK, true, "删除成功", nil)
}
