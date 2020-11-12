// +build test public api tag access perm role user type doc chapter article expire config article dashboard

package tests

import (
	"flag"
	"github.com/bxcodec/faker/v3"
	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/seeder"
	"github.com/snowlyg/blog/tests/mock"
	"github.com/snowlyg/blog/web_server"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/snowlyg/blog/libs"
)

var (
	ConfigPath = flag.String("tc", "", "配置路径")
	app        *iris.Application
	token      string
)

//单元测试基境
func TestMain(m *testing.M) {
	flag.Parse()
	libs.InitConfig(*ConfigPath)
	s := web_server.NewServer(nil) // 初始化app
	s.NewApp()
	app = s.App
	seeder.Run()

	exitCode := m.Run()

	models.DropTables() // 删除测试数据表，保持测试环境
	os.Exit(exitCode)
}

func getHttpexpect(t *testing.T) *httpexpect.Expect {
	return httptest.New(t, app, httptest.Configuration{Debug: true, URL: "http://app.irisadminapi.com/v1/admin/"})
}
func getPublicHttpexpect(t *testing.T) *httpexpect.Expect {
	return httptest.New(t, app, httptest.Configuration{Debug: true, URL: "http://app.irisadminapi.com/v1/"})
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

// 单元测试 getOne 方法
func getOnePublic(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getPublicHttpexpect(t)
	e.GET(url).
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

// 单元测试 getPublicMore 方法
func getPublicMore(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getPublicHttpexpect(t)
	e.GET(url).
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

func CreateArticle(status string) (*models.Article, error) {
	mock.CustomGenerator()
	m := mock.Article{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	if len(status) == 0 {
		status = m.Status
	}
	article := &models.Article{
		Title:        m.Title,
		ContentShort: m.ContentShort,
		Author:       m.Author,
		ImageUri:     m.ImageUri,
		SourceUri:    m.SourceUri,
		IsOriginal:   m.IsOriginal,
		Content:      m.ContentShort,
		Status:       status,
		DisplayTime:  time.Now(),
		Like:         m.Like,
		Read:         m.Read,
		Ips:          m.Ips,
	}
	err = article.CreateArticle()
	if err != nil {
		return article, err
	}
	return article, nil
}
func CreateChapter(status string) (*models.Chapter, error) {
	mock.CustomGenerator()
	m := mock.Chapter{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}

	if len(status) == 0 {
		status = m.Status
	}

	chapter := &models.Chapter{
		Title:        m.Title,
		ContentShort: m.ContentShort,
		Author:       m.Author,
		ImageUri:     m.ImageUri,
		SourceUri:    m.SourceUri,
		IsOriginal:   m.IsOriginal,
		Content:      m.ContentShort,
		Status:       status,
		DisplayTime:  time.Now(),
		Like:         m.Like,
		Read:         m.Read,
		Ips:          m.Ips,
		Sort:         m.Sort,
	}
	err = chapter.CreateChapter()
	if err != nil {
		return chapter, err
	}
	return chapter, nil
}

func CreateConfig() (*models.Config, error) {
	m := mock.Config{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	config := &models.Config{
		Name:  m.Name,
		Value: m.Value,
	}
	err = config.CreateConfig()
	if err != nil {
		return config, err
	}
	return config, nil
}

func CreateDoc() (*models.Doc, error) {
	m := mock.Doc{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	doc := &models.Doc{
		Name: m.Name,
	}
	err = doc.CreateDoc()
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func CreateType() (*models.Type, error) {
	m := mock.Type{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	tt := &models.Type{
		Name: m.Name,
	}
	err = tt.CreateType()
	if err != nil {
		return tt, err
	}
	return tt, nil
}

func CreateTag() (*models.Tag, error) {
	m := mock.Tag{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	tag := &models.Tag{
		Name: m.Name,
	}
	err = tag.CreateTag()
	if err != nil {
		return tag, err
	}
	return tag, nil
}

func CreatePermission() (*models.Permission, error) {
	m := mock.Permission{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	perm := &models.Permission{
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Description: m.Description,
		Act:         m.Act,
	}
	err = perm.CreatePermission()
	if err != nil {
		return perm, err
	}

	return perm, nil

}

func CreateRole() (*models.Role, error) {
	m := mock.Role{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	role := &models.Role{
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Description: m.Description,
	}
	err = role.CreateRole()
	if err != nil {
		return role, err
	}

	return role, nil
}

func CreateUser() (*models.User, error) {
	r, err := CreateRole()
	if err != nil {
		return nil, err
	}
	m := mock.User{}
	err = faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username: m.Username,
		Password: m.Password,
		Name:     m.Name,
		RoleIds:  []uint{r.ID},
	}
	err = user.CreateUser()
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
