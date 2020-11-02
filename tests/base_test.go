// +build test api tag access perm role user type doc chapter article

package tests

import (
	"flag"
	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/seeder"
	"github.com/snowlyg/blog/web_server"
	"net/http"
	"os"
	"testing"

	"github.com/snowlyg/blog/libs"
)

var (
	app   *iris.Application
	token string
)

//单元测试基境
func TestMain(m *testing.M) {
	s := web_server.NewServer(nil) // 初始化app
	s.NewApp()
	app = s.App
	seeder.Run()

	flag.Parse()
	exitCode := m.Run()

	models.DropTables() // 删除测试数据表，保持测试环境
	os.Exit(exitCode)
}

func getHttpexpect(t *testing.T) *httpexpect.Expect {
	return httptest.New(t, app, httptest.Configuration{Debug: true, URL: "http://app.irisadminapi.com/v1/admin/"})
}

// 单元测试 login 方法
func login(t *testing.T, Object interface{}, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.POST("login").WithJSON(Object).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)

	return
}

// 单元测试 create 方法
func create(t *testing.T, url string, Object interface{}, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	ob := e.POST(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()
	ob.Value("code").Equal(Code)
	ob.Value("message").Equal(Msg)

	return
}

// 单元测试 update 方法
func update(t *testing.T, url string, Object interface{}, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	ob := e.PUT(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()
	ob.Value("code").Equal(Code)
	ob.Value("message").Equal(Msg)

	return
}

// 单元测试 getOne 方法
func getOne(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)
	return
}

// 单元测试 getOnAuth 方法
func getOnAuth(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.GET(url).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)
	return
}

// 单元测试 bImport 方法
func bImport(t *testing.T, url string, StatusCode int, Code int, Msg string, _ map[string]interface{}) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.POST(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		WithMultipart().
		WithFile("file", "permissions.xlsx").
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)

	return
}

// 单元测试 getMore 方法
func getMore(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)

	return
}

// 单元测试 delete 方法
func delete(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.DELETE(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)
	return
}

func CreateDoc(name string) (*models.Doc, error) {
	doc := &models.Doc{
		Name: name,
	}
	err := doc.CreateDoc()
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func CreateType(name string) (*models.Type, error) {
	tt := &models.Type{
		Name: name,
	}
	err := tt.CreateType()
	if err != nil {
		return tt, err
	}
	return tt, nil
}

func CreateTag(name string) (*models.Tag, error) {
	tag := &models.Tag{
		Name: name,
	}
	err := tag.CreateTag()
	if err != nil {
		return tag, err
	}
	return tag, nil
}

func CreatePermission(name, disName, dec string) (*models.Permission, error) {
	perm := &models.Permission{
		Name:        name,
		DisplayName: disName,
		Description: dec,
	}
	err := perm.CreatePermission()
	if err != nil {
		return perm, err
	}

	return perm, nil

}

func CreateRole(name, disName, dec string) (*models.Role, error) {
	role := &models.Role{
		Name:        name,
		DisplayName: disName,
		Description: dec,
	}
	err := role.CreateRole()
	if err != nil {
		return role, err
	}

	return role, nil
}

func CreateUser() (*models.User, error) {
	user := &models.User{
		Username: "TUsername",
		Password: "TPassword",
		Name:     "TName",
		RoleIds:  []uint{},
	}
	err := user.CreateUser()
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetOauthToken(e *httpexpect.Expect) string {
	if len(token) > 0 {
		return token
	}

	oj := map[string]string{
		"username": libs.Config.Admin.UserName,
		"password": libs.Config.Admin.Pwd,
	}
	r := e.POST("login").WithJSON(oj).
		Expect().
		Status(http.StatusOK).JSON().Object()

	token = r.Value("data").Object().Value("token").String().Raw()

	return token
}
