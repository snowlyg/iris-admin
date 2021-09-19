package test

import (
	"fmt"
	"testing"

	"github.com/snowlyg/helper/tests"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/modules/v1/perm"
	"github.com/snowlyg/iris-admin/server/database"
)

var (
	loginUrl  = "/api/v1/auth/login"
	logoutUrl = "/api/v1/users/logout"
	url       = "/api/v1/perms"
	data      = map[string]interface{}{
		"name":        "test_route_name",
		"displayName": "测试描述信息",
		"description": "测试描述信息",
		"act":         "GET",
	}
)

type PageParam struct {
	Message  string
	Code     int
	PageSize int
	Page     int
	PageLen  int
}

func TestList(t *testing.T) {
	client := TestServer.GetTestLogin(t, loginUrl, nil)
	defer client.Logout(logoutUrl, nil)

	pageParams := getPageParams()
	for _, pageParam := range pageParams {
		t.Run(fmt.Sprintf("路由权限测试，第%d页", pageParam.Page), func(t *testing.T) {
			items, err := getPerms(pageParam)
			if err != nil {
				t.Fatalf("获取路由权限错误")
			}
			pageKeys := tests.Responses{
				{Key: "code", Value: pageParam.Code},
				{Key: "message", Value: pageParam.Message},
				{Key: "data", Value: tests.Responses{
					{Key: "pageSize", Value: pageParam.PageSize},
					{Key: "page", Value: pageParam.Page},
					{Key: "items", Value: items},
					{Key: "total", Value: len(g.PermRoutes)},
				}},
			}
			requestParams := map[string]interface{}{"page": pageParam.Page, "pageSize": pageParam.PageSize}
			client.GET(url, pageKeys, requestParams)
		})
	}

}

func TestCreate(t *testing.T) {
	client := TestServer.GetTestLogin(t, loginUrl, nil)
	defer client.Logout(logoutUrl, nil)
	userId := Create(client, data)
	if userId == 0 {
		t.Fatalf("测试添加用户失败 id=%d", userId)
	}
	defer Delete(client, userId)
}

func TestUpdate(t *testing.T) {
	client := TestServer.GetTestLogin(t, loginUrl, nil)
	defer client.Logout(logoutUrl, nil)
	userId := Create(client, data)
	if userId == 0 {
		t.Fatalf("测试添加用户失败 id=%d", userId)
	}
	defer Delete(client, userId)

	update := map[string]interface{}{
		"name":        "update_test_route_name",
		"displayName": "更新测试描述信息",
		"description": "更新测试描述信息",
		"act":         "POST",
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
			{Key: "act", Value: data["act"].(string)},
			{Key: "updatedAt", Value: "", Type: "notempty"},
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

func getPerms(pageParam PageParam) ([]tests.Responses, error) {
	l := pageParam.PageLen
	routes := make([]tests.Responses, 0, l)
	req := perm.ReqPaginate{
		Paginate: g.Paginate{
			Page:     pageParam.Page,
			PageSize: pageParam.PageSize,
		},
	}
	perms, err := perm.Paginate(database.Instance(), req)
	if err != nil {
		return routes, err
	}
	for _, route := range perms["items"].([]*perm.Response) {
		perm := tests.Responses{
			{Key: "id", Value: route.Id},
			{Key: "name", Value: route.Name},
			{Key: "displayName", Value: route.DisplayName},
			{Key: "description", Value: route.Description},
			{Key: "act", Value: route.Act},
			{Key: "updatedAt", Value: route.UpdatedAt},
			{Key: "createdAt", Value: route.CreatedAt},
		}
		routes = append(routes, perm)
		l--
		if l == 0 {
			break
		}
	}

	return routes, err
}

func getPageParams() []PageParam {
	pageSize := 10
	size := len(g.PermRoutes) / pageSize
	other := len(g.PermRoutes) % pageSize
	if other > 0 {
		size++
	}
	pageParams := make([]PageParam, 0, size)
	for i := 1; i <= size; i++ {
		pageLen := pageSize
		if other > 0 && i == size {
			pageLen = other
		}
		pageParam := PageParam{
			Message:  "请求成功",
			Code:     2000,
			PageSize: pageSize,
			PageLen:  pageLen,
			Page:     i,
		}
		pageParams = append(pageParams, pageParam)
	}
	return pageParams
}
