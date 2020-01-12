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
	Sc = iris.TOML("./config/conf.tml")
	rc = getSysConf()

	// 初始化app
	app = NewApp(rc)

	flag.Parse()
	exitCode := m.Run()

	// 删除测试数据表，保持测试环境
	models.Db.DropTable("users", "roles", "permissions", "oauth_tokens")

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

	ob := e.POST(url).WithHeader("Authorization", "Bearer "+GetOauthToken()).WithJSON(Object).
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

	ob := e.PUT(url).WithHeader("Authorization", "Bearer "+GetOauthToken()).WithJSON(Object).
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
	if Data != nil {
		e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken()).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg, Data)
	} else {
		e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken()).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg)
	}

	return
}

// 单元测试 getMore 方法
func getMore(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})

	if Data != nil {
		e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken()).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg, Data)
	} else {
		e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken()).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg)
	}

	return
}

// 单元测试 delete 方法
func delete(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})

	e.DELETE(url).WithHeader("Authorization", "Bearer "+GetOauthToken()).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)

	return
}

func CreateRole(name, disName, dec string) *models.Role {
	rr := &models.RoleRequest{
		Name:        name,
		DisplayName: disName,
		Description: dec,
	}

	role := models.GetRoleByName("Tname")
	if role.ID == 0 {
		return models.CreateRole(rr, []uint{})
	} else {
		return role
	}
}

func CreateUser() *models.User {
	rr := &models.UserRequest{
		Username: "TUsername",
		Password: "TPassword",
		Name:     "TName",
		RoleIds:  []uint{},
	}

	user := models.GetUserByUserName("TUsername")
	if user.ID == 0 {
		return models.CreateUser(rr)
	} else {
		return user
	}
}

func GetOauthToken() string {
	ot, b, msg := models.CheckLogin(rc.TestData.UserName, rc.TestData.Pwd)
	if b {
		return ot.Token
	} else {
		fmt.Println(fmt.Sprintf("GetOauthToken Error : %v", msg))
		return ""
	}
}
