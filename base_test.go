package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"IrisAdminApi/models"
	"IrisAdminApi/transformer"
	"github.com/gavv/httpexpect"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

var (
	app *iris.Application // iris.Applications
	rc  *transformer.Conf
)

//单元测试基境
func TestMain(m *testing.M) {

	// 设置静态资源
	Sc = iris.TOML("./config/test.tml")
	rc = getSysConf()

	// 初始化app
	app = NewApp(rc)

	flag.Parse()
	exitCode := m.Run()

	// 删除测试数据表，保持测试环境
	models.Db.DropTable("users", "roles", "permissions", &models.OauthToken{})

	os.Exit(exitCode)

}

// 单元测试 login 方法
func login(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	if Data != nil {
		e.POST(url).WithJSON(Object).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg, Data)
	} else {
		e.POST(url).WithJSON(Object).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg)
	}

	return
}

// 单元测试 create 方法
func create(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	at := GetLoginToken()

	ob := e.POST(url).WithHeader("Authorization", "Bearer "+at.Token).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()

	ob.Value("status").Equal(Status)
	ob.Value("msg").Equal(Msg)

	for k, v := range Data {
		ob.Value("data").Object().Value(k).Equal(v)
	}

	return
}

// 单元测试 update 方法
func update(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	at := GetLoginToken()

	ob := e.PUT(url).WithHeader("Authorization", "Bearer "+at.Token).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()

	ob.Value("status").Equal(Status)
	ob.Value("msg").Equal(Msg)

	for k, v := range Data {
		ob.Value("data").Object().Value(k).Equal(v)
	}

	return
}

// 单元测试 getOne 方法
func getOne(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	at := GetLoginToken()
	if Data != nil {
		e.GET(url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg, Data)
	} else {
		e.GET(url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg)
	}

	return
}

// 单元测试 getMore 方法
func getMore(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	at := GetLoginToken()
	if Data != nil {
		e.GET(url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg, Data)
	} else {
		e.GET(url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg)
	}

	return
}

// 单元测试 delete 方法
func delete(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	at := GetLoginToken()

	e.DELETE(url).WithHeader("Authorization", "Bearer "+at.Token).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)

	return
}

/**
*登陆用户
*@return   Token 返回登陆后的token
 */
func GetLoginToken() models.Token {
	response, status, msg := models.CheckLogin(
		rc.TestData.UserName,
		rc.TestData.Pwd,
	)

	// 打印错误信息
	if !status {
		fmt.Println(msg)
	}

	return response
}

func CreateRole() *models.Role {
	rr := &models.RoleRequest{
		Name:        "name",
		DisplayName: "DisplayName",
		Description: "DisplayName",
	}

	return models.CreateRole(rr, []uint{})
}

func CreateUser() *models.User {
	rr := &models.UserRequest{
		Username: "Username",
		Password: "Password",
		Name:     "Name",
		RoleIds:  []uint{},
	}

	return models.CreateUser(rr)
}
