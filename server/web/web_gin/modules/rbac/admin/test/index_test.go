package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/helper/tests"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	v1 "github.com/snowlyg/iris-admin/server/web/web_gin/modules/rbac"
	"github.com/snowlyg/iris-admin/server/web/web_gin/response"
)

var (
	loginUrl  = "/api/v1/public/admin/login"
	logoutUrl = "/api/v1/public/logout"
	url       = "/api/v1/admin"
)

func TestList(t *testing.T) {
	if TestServer == nil {
		t.Error("test server is nil")
		return
	}

	client := TestServer.GetTestLogin(t, loginUrl, v1.LoginResponse)
	if client != nil {
		defer client.Logout(logoutUrl, v1.LogoutResponse)
	} else {
		return
	}
	pageKeys := tests.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
		{Key: "data", Value: tests.Responses{
			{Key: "pageSize", Value: 10},
			{Key: "page", Value: 1},
			{Key: "items", Value: []tests.Responses{
				{
					{Key: "id", Value: 1, Type: "ge"},
					{Key: "name", Value: "超级管理员"},
					{Key: "username", Value: "admin"},
					{Key: "intro", Value: "超级管理员"},
					{Key: "avatar", Value: str.Join("http://", web_gin.CONFIG.System.Addr, web_gin.CONFIG.System.StaticPrefix, "/images/avatar.jpg")},
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
	if TestServer == nil {
		t.Error("test server is nil")
		return
	}

	client := TestServer.GetTestLogin(t, loginUrl, v1.LoginResponse)
	if client != nil {
		defer client.Logout(logoutUrl, v1.LogoutResponse)
	} else {
		return
	}
	data := map[string]interface{}{
		"name":     "测试名称",
		"username": "create_test_username",
		"intro":    "测试描述信息",
		"avatar":   "",
		"password": "123456",
	}
	id := Create(client, data)
	if id == 0 {
		t.Fatalf("测试添加用户失败 id=%d", id)
	}
	defer Delete(client, id)
}

func TestUpdate(t *testing.T) {
	if TestServer == nil {
		t.Error("test server is nil")
		return
	}

	client := TestServer.GetTestLogin(t, loginUrl, v1.LoginResponse)
	if client != nil {
		defer client.Logout(logoutUrl, v1.LogoutResponse)
	} else {
		return
	}
	data := map[string]interface{}{
		"name":     "测试名称",
		"username": "update_test_username",
		"intro":    "测试描述信息",
		"avatar":   "",
		"password": "123456",
	}
	id := Create(client, data)
	if id == 0 {
		t.Fatalf("测试添加用户失败 id=%d", id)
	}
	defer Delete(client, id)

	update := map[string]interface{}{
		"name":     "更新测试名称",
		"username": "update_test_username",
		"intro":    "更新测试描述信息",
		"avatar":   "",
		"password": "123456",
	}

	pageKeys := tests.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
	}
	client.POST(fmt.Sprintf("%s/%d", url, id), pageKeys, update)
}

func TestGetById(t *testing.T) {
	if TestServer == nil {
		t.Error("test server is nil")
		return
	}

	client := TestServer.GetTestLogin(t, loginUrl, v1.LoginResponse)
	if client != nil {
		defer client.Logout(logoutUrl, v1.LogoutResponse)
	} else {
		return
	}
	data := map[string]interface{}{
		"name":     "测试名称",
		"username": "getbyid_test_username",
		"intro":    "测试描述信息",
		"avatar":   "",
		"password": "123456",
	}
	id := Create(client, data)
	if id == 0 {
		t.Fatalf("测试添加用户失败 id=%d", id)
	}
	defer Delete(client, id)

	pageKeys := tests.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
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
	client.GET(fmt.Sprintf("%s/%d", url, id), pageKeys)
}

func Create(client *tests.Client, data map[string]interface{}) uint {
	pageKeys := tests.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
		{Key: "data", Value: tests.Responses{
			{Key: "id", Value: 1, Type: "ge"},
		},
		},
	}
	return client.POST(url, pageKeys, data).GetId()
}

func Delete(client *tests.Client, id uint) {
	pageKeys := tests.Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: response.ResponseOkMessage},
	}
	client.DELETE(fmt.Sprintf("%s/%d", url, id), pageKeys)
}
