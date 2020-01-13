package main

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"IrisAdminApi/models"
	"IrisAdminApi/transformer"
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
func getOne(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	if Data != nil {
		e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg, Data)
	} else {
		e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg)
	}

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
