package test

import (
	"fmt"
	"testing"

	"github.com/snowlyg/helper/tests"
)

var (
	loginUrl = "/api/v1/auth/login" // 登录URL
	url      = "/api/v1/roles"      // url
)

func TestList(t *testing.T) {
	if TestServer == nil {
		t.Errorf("TestServer is nil")
	}
	TestClient = TestServer.GetTestLogin(t, loginUrl, nil)
	if TestClient == nil {
		return
	}
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
	TestClient.GET(url, pageKeys, tests.RequestParams)
}

func TestCreate(t *testing.T) {
	if TestServer == nil {
		t.Errorf("TestServer is nil")
	}
	TestClient = TestServer.GetTestLogin(t, loginUrl, nil)
	if TestClient == nil {
		return
	}
	data := map[string]interface{}{
		"name":        "test_display_name",
		"displayName": "测试名称",
		"description": "测试描述信息",
	}
	id := Create(TestClient, data)
	if id == 0 {
		t.Fatalf("测试添加用户失败 id=%d", id)
	}
	defer Delete(TestClient, id)
}

func TestUpdate(t *testing.T) {
	if TestServer == nil {
		t.Errorf("TestServer is nil")
	}
	TestClient = TestServer.GetTestLogin(t, loginUrl, nil)
	if TestClient == nil {
		return
	}
	data := map[string]interface{}{
		"name":        "update_test_display_name",
		"displayName": "测试名称",
		"description": "测试描述信息",
	}
	id := Create(TestClient, data)
	if id == 0 {
		t.Fatalf("测试添加用户失败 id=%d", id)
	}
	defer Delete(TestClient, id)

	update := map[string]interface{}{
		"name":        "update_test_udisplay_name",
		"displayName": "更新测试名称",
		"description": "更新测试描述信息",
	}

	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
	}
	TestClient.POST(fmt.Sprintf("%s/%d", url, id), pageKeys, update)
}

func TestGetById(t *testing.T) {
	if TestServer == nil {
		t.Errorf("TestServer is nil")
	}
	TestClient = TestServer.GetTestLogin(t, loginUrl, nil)
	if TestClient == nil {
		return
	}
	data := map[string]interface{}{
		"name":        "getbyid_test_display_name",
		"displayName": "更新测试名称",
		"description": "测试描述信息",
	}
	id := Create(TestClient, data)
	if id == 0 {
		t.Fatalf("测试添加失败 id=%d", id)
	}
	defer Delete(TestClient, id)

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
	TestClient.GET(fmt.Sprintf("%s/%d", url, id), pageKeys)
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
