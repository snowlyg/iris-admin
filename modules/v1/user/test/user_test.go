package test

import (
	"fmt"
	"testing"

	"github.com/snowlyg/helper/tests"
)

func TestList(t *testing.T) {
	client := TestServer.GetTestLogin(t, "/api/v1/auth/login", nil)
	defer client.Logout("/api/v1/users/logout", nil)
	url := "/api/v1/users"
	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
		{Key: "data", Value: tests.Responses{
			{Key: "pageSize", Value: 10},
			{Key: "page", Value: 1},
			{Key: "items", Value: []tests.Responses{
				{
					{Key: "id", Value: 1, Type: "ge"},
					{Key: "name", Value: "超级管理员"},
					{Key: "username", Value: "admin"},
					{Key: "intro", Value: "超级管理员"},
					{Key: "avatar", Value: "/static/images/avatar.jpg"},
					{Key: "roles", Value: []string{"超级管理员"}},
					{Key: "updatedAt", Value: "", Type: "notempty"},
					{Key: "createdAt", Value: "", Type: "notempty"},
				},
			}},
			{Key: "total", Value: 0, Type: "ge"},
		}},
	}
	client.GET(url, pageKeys, tests.RequestParams)
}

func TestCreate(t *testing.T) {
	client := TestServer.GetTestLogin(t, "/api/v1/auth/login", nil)
	defer client.Logout("/api/v1/users/logout", nil)
	data := map[string]interface{}{
		"name":     "测试名称",
		"username": "test_username",
		"intro":    "测试描述信息",
		"avatar":   "",
		"password": "123456",
	}
	userId := CreateUser(client, data)
	if userId == 0 {
		t.Fatalf("测试添加用户失败 id=%d", userId)
	}
	defer DeleteUser(client, userId)
}

func TestUpdate(t *testing.T) {
	client := TestServer.GetTestLogin(t, "/api/v1/auth/login", nil)
	defer client.Logout("/api/v1/users/logout", nil)
	data := map[string]interface{}{
		"name":     "测试名称",
		"username": "test_username",
		"intro":    "测试描述信息",
		"avatar":   "",
		"password": "123456",
	}
	userId := CreateUser(client, data)
	if userId == 0 {
		t.Fatalf("测试添加用户失败 id=%d", userId)
	}
	defer DeleteUser(client, userId)

	update := map[string]interface{}{
		"name":     "更新测试名称",
		"username": "update_test_username",
		"intro":    "更新测试描述信息",
		"avatar":   "",
		"password": "123456",
	}

	url := fmt.Sprintf("/api/v1/users/%d", userId)
	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
	}
	client.POST(url, pageKeys, update)
}

func TestGetById(t *testing.T) {
	client := TestServer.GetTestLogin(t, "/api/v1/auth/login", nil)
	defer client.Logout("/api/v1/users/logout", nil)
	data := map[string]interface{}{
		"name":     "测试名称",
		"username": "test_username",
		"intro":    "测试描述信息",
		"avatar":   "",
		"password": "123456",
	}
	userId := CreateUser(client, data)
	if userId == 0 {
		t.Fatalf("测试添加用户失败 id=%d", userId)
	}
	defer DeleteUser(client, userId)
	url := fmt.Sprintf("/api/v1/users/%d", userId)
	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
		{Key: "data", Value: tests.Responses{
			{Key: "id", Value: 1, Type: "ge"},
			{Key: "name", Value: data["name"].(string)},
			{Key: "username", Value: data["username"].(string)},
			{Key: "intro", Value: data["intro"].(string)},
			{Key: "avatar", Value: data["avatar"].(string)},
			{Key: "updatedAt", Value: "", Type: "notempty"},
			{Key: "createdAt", Value: "", Type: "notempty"},
			{Key: "createdAt", Value: "", Type: "notempty"},
			{Key: "roles", Value: []string{}, Type: "null"},
		},
		},
	}
	client.GET(url, pageKeys)
}

func CreateUser(client *tests.Client, data map[string]interface{}) uint {
	url := "/api/v1/users"
	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
		{Key: "data", Value: tests.Responses{
			{Key: "id", Value: 1, Type: "ge"},
		},
		},
	}
	return client.POST(url, pageKeys, data).GetId()
}

func DeleteUser(client *tests.Client, id uint) {
	url := fmt.Sprintf("/api/v1/users/%d", id)
	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
	}
	client.DELETE(url, pageKeys)
}
