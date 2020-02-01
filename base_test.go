package main

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"IrisAdminApi/config"
	"IrisAdminApi/database"
	"IrisAdminApi/models"
	"IrisAdminApi/routepath"
	"IrisAdminApi/transformer"
	"IrisAdminApi/validates"
	"github.com/gavv/httpexpect"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

const baseUrl = "/v1/admin/"
const loginUrl = baseUrl + "login"

var (
	app   *iris.Application // iris.Applications
	rc    *transformer.Conf
	token string
)

//单元测试基境
func TestMain(m *testing.M) {

	rc = config.GetTfConf()
	app = NewApp(rc) // 初始化app
	routes := routepath.GetRoutes(app.APIBuilder.GetRoutes())
	models.CreateSystemData(rc, routes)

	flag.Parse()
	exitCode := m.Run()

	database.DropTables() // 删除测试数据表，保持测试环境
	os.Exit(exitCode)
}

// 单元测试 login 方法
func login(t *testing.T, Object interface{}, StatusCode int, Status bool, Msg string) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	e.POST(loginUrl).WithJSON(Object).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)

	return
}

// 单元测试 create 方法
func create(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	ob := e.POST(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()
	ob.Value("status").Equal(Status)
	ob.Value("msg").Equal(Msg)

	return
}

// 单元测试 update 方法
func update(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	ob := e.PUT(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()
	ob.Value("status").Equal(Status)
	ob.Value("msg").Equal(Msg)

	return
}

// 单元测试 getOne 方法
func getOne(t *testing.T, url string, StatusCode int, Status bool, Msg string) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)
	return
}

// 单元测试 getOnAuth 方法
func getOnAuth(t *testing.T, url string, StatusCode int, Status bool, Msg string) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	e.GET(url).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)
	return
}

// 单元测试 bImport 方法
func bImport(t *testing.T, url string, StatusCode int, Status bool, Msg string, _ map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	e.POST(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		WithMultipart().
		WithFile("file", "permissions.xlsx").
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)

	return
}

// 单元测试 getMore 方法
func getMore(t *testing.T, url string, StatusCode int, Status bool, Msg string) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)

	return
}

// 单元测试 delete 方法
func delete(t *testing.T, url string, StatusCode int, Status bool, Msg string) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	e.DELETE(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)
	return
}

func CreateRole(name, disName, dec string) *models.Role {
	rr := &validates.RoleRequest{
		Name:        name,
		DisplayName: disName,
		Description: dec,
	}
	role := models.NewRoleByStruct(rr)
	role.GetRoleByName()
	if role.ID == 0 {
		role.CreateRole([]uint{})
	}
	return role

}

func CreateUser() *models.User {
	rr := &validates.CreateUpdateUserRequest{
		Username: "TUsername",
		Password: "TPassword",
		Name:     "TName",
		RoleIds:  []uint{},
	}

	user := models.NewUser(0, "TUsername")
	if user.ID == 0 {
		user.CreateUser(rr)
		return user
	} else {
		return user
	}
}

func GetOauthToken(e *httpexpect.Expect) string {
	if len(token) > 0 {
		return token
	}

	oj := map[string]string{
		"username": rc.TestData.UserName,
		"password": rc.TestData.Pwd,
	}
	r := e.POST(loginUrl).WithJSON(oj).
		Expect().
		Status(http.StatusOK).JSON().Object()

	token = r.Value("data").Object().Value("access_token").String().Raw()

	return token
}
