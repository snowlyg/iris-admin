package test

import (
	"fmt"
	"testing"

	"github.com/snowlyg/helper/tests"
)

var (
	loginUrl  = "/api/v1/auth/login"   // 登录URL
	logoutUrl = "/api/v1/users/logout" // 登出 URL
	url       = "/api/v1/roles"        // url
)

func TestList(t *testing.T) {
	client := TestServer.GetTestLogin(t, loginUrl, nil)
	defer client.Logout(logoutUrl, nil)
	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
		{Key: "data", Value: tests.Responses{
			{Key: "pageSize", Value: 10},
			{Key: "page", Value: 1},
			{Key: "items", Value: []tests.Responses{
				{
					{Key: "id", Value: 1, Type: "ge"},
					{Key: "name", Value: "SuperAdmin"},
					{Key: "displayName", Value: "超级管理员"},
					{Key: "description", Value: "超级管理员"},
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
	client := TestServer.GetTestLogin(t, loginUrl, nil)
	defer client.Logout(logoutUrl, nil)
	data := map[string]interface{}{
		"name":        "测试名称",
		"displayName": "test_display_name",
		"description": "测试描述信息",
	}
	userId := Create(client, data)
	if userId == 0 {
		t.Fatalf("测试添加用户失败 id=%d", userId)
	}
	defer Delete(client, userId)
}

func TestUpdate(t *testing.T) {
	client := TestServer.GetTestLogin(t, loginUrl, nil)
	defer client.Logout(logoutUrl, nil)
	data := map[string]interface{}{
		"name":        "测试名称",
		"displayName": "update_test_display_name",
		"description": "测试描述信息",
	}
	userId := Create(client, data)
	if userId == 0 {
		t.Fatalf("测试添加用户失败 id=%d", userId)
	}
	defer Delete(client, userId)

	update := map[string]interface{}{
		"name":        "更新测试名称",
		"displayName": "update_test_udisplay_name",
		"description": "更新测试描述信息",
	}

	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
	}
	client.POST(fmt.Sprintf("%s/%d", url, userId), pageKeys, update)
}

func TestGetById(t *testing.T) {
	client := TestServer.GetTestLogin(t, loginUrl, nil)
	defer client.Logout(logoutUrl, nil)
	data := map[string]interface{}{
		"name":        "测试名称",
		"displayName": "getbyid_test_display_name",
		"description": "测试描述信息",
	}
	userId := Create(client, data)
	if userId == 0 {
		t.Fatalf("测试添加用户失败 id=%d", userId)
	}
	defer Delete(client, userId)

	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
		{Key: "data", Value: tests.Responses{
			{Key: "id", Value: 1, Type: "ge"},
			{Key: "name", Value: data["name"].(string)},
			{Key: "displayName", Value: data["displayName"].(string)},
			{Key: "description", Value: data["description"].(string)},
			{Key: "updatedAt", Value: "", Type: "notempty"},
			{Key: "createdAt", Value: "", Type: "notempty"},
			{Key: "createdAt", Value: "", Type: "notempty"},
		},
		},
	}
	client.GET(fmt.Sprintf("%s/%d", url, userId), pageKeys)
}

func Create(client *tests.Client, data map[string]interface{}) uint {
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

func Delete(client *tests.Client, id uint) {
	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
	}
	client.DELETE(fmt.Sprintf("%s/%d", url, id), pageKeys)
}
