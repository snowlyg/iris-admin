// +build test public api tag access perm role user type doc chapter article expire config article dashboard

package tests

import (
	"github.com/bxcodec/faker/v3"
	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/snowlyg/blog/app"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/seeder"
	"github.com/snowlyg/blog/tests/mock"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/snowlyg/blog/libs"
)

var (
	application  *iris.Application
	token        string
	ArticleCount int64
	DocCount     int64
	TypeCount    int64
	TagCount     int64
	ChapterCount int64
	UserCount    int64 = 1
	RoleCount    int64 = 1
	PermCount    int64 = 53
	ConfigCount  int64 = 2
)

//基境
func TestMain(m *testing.M) {

	libs.InitConfig("")
	easygorm.Init(&easygorm.Config{
		GormConfig: &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "iris_", // 表名前缀，`User` 的表名应该是 `t_users`
				SingularTable: false,   // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
			},
		},
		Adapter:         "sqlite3", // 类型
		Name:            "blog_test",
		Username:        "root",                                       // 用户名
		Pwd:             "123456",                                     // 密码
		Host:            "127.0.0.1",                                  // 地址
		Port:            3306,                                         // 端口
		CasbinModelPath: filepath.Join(libs.CWD(), "rbac_model.conf"), // casbin 模型规则路径
		Debug:           true,
		TablePrefix:     "iris",
		Models: []interface{}{
			&models.User{},
			&models.Role{},
			&models.Permission{},
			&models.Article{},
			&models.Config{},
			&models.Tag{},
			&models.Type{},
			&models.Doc{},
			&models.Chapter{},
			&models.ChapterIp{},
			&models.ArticleIp{},
		},
	})
	s := app.NewServer() // 初始化app
	application = s.App
	seeder.Run()

	exitCode := m.Run()

	models.DropTables("iris_") // 删除测试数据表，保持测试环境
	os.Exit(exitCode)
}

func getHttpexpect(t *testing.T) *httpexpect.Expect {
	return httptest.New(t, application, httptest.Configuration{Debug: true, URL: "http://app.irisadminapi.com/v1/admin/"})
}
func getPublicHttpexpect(t *testing.T) *httpexpect.Expect {
	return httptest.New(t, application, httptest.Configuration{Debug: true, URL: "http://app.irisadminapi.com/v1/"})
}

//  login 方法
func login(t *testing.T, Object interface{}, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.POST("login").WithJSON(Object).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)

	return
}

//  create 方法
func create(t *testing.T, url string, Object interface{}, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	ob := e.POST(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()
	ob.Value("code").Equal(Code)
	ob.Value("message").Equal(Msg)

	return
}

//  update 方法
func update(t *testing.T, url string, Object interface{}, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	ob := e.PUT(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()
	ob.Value("code").Equal(Code)
	ob.Value("message").Equal(Msg)

	return
}

//  getOne 方法
func getOne(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)
	return
}

//  getOne 方法
func getOnePublic(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getPublicHttpexpect(t)
	e.GET(url).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)
	return
}

//  getNoAuth 方法
func getNoAuth(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.GET(url).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)
	return
}

//  bImport 方法
func bImport(t *testing.T, url string, StatusCode int, Code int, Msg string, _ map[string]interface{}) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.POST(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		WithMultipart().
		WithFile("file", "permissions.xlsx").
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)

	return
}

type More struct {
	Id    uint
	Limit int
	Page  int64
	Total int64
	Field []interface{}
}

//  getNoReturn 方法
func getNoReturn(t *testing.T, url string, StatusCode int, obj interface{}) {
	e := getHttpexpect(t)
	o := e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		WithQueryObject(obj).
		Expect().Status(StatusCode).
		JSON().Object()

	o.Value("code").Equal(200)
	o.Value("message").Equal("")
	o.Value("data").Null()
}

//  getMore 方法
func getMore(t *testing.T, url string, StatusCode int, obj interface{}, m *More) {
	e := getHttpexpect(t)

	o := e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		WithQueryObject(obj).
		Expect().Status(StatusCode).
		JSON().Object()

	o.Value("code").Equal(200)
	o.Value("message").Equal("操作成功")
	items := o.Value("data").Object().Value("items").Array()
	items.Length().Equal(m.Page)
	items.First().Object().Value("id").Equal(m.Id)
	items.First().Object().Keys().Contains(m.Field...)
	o.Value("data").Object().Value("limit").Equal(m.Limit)
	o.Value("data").Object().Value("total").Equal(m.Total)
}

//  getAll 方法
func getAll(t *testing.T, url string, StatusCode int, obj interface{}, total int64) {
	e := getHttpexpect(t)
	o := e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		WithQueryObject(obj).
		Expect().Status(StatusCode).
		JSON().Object()

	o.Value("code").Equal(200)
	o.Value("message").Equal("操作成功")
	o.Value("data").Array().Length().Equal(total)
}

//  getData 方法
func getData(t *testing.T, url string, StatusCode int, obj interface{}, keys []interface{}) {
	e := getHttpexpect(t)
	o := e.GET(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		WithQueryObject(obj).
		Expect().Status(StatusCode).
		JSON().Object()

	o.Value("code").Equal(200)
	o.Value("message").Equal("操作成功")
	o.Value("data").Object().Keys().Contains(keys...)
}

//  getPublicMore 方法
func getPublicMore(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getPublicHttpexpect(t)
	e.GET(url).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)

	return
}

//  delete 方法
func delete(t *testing.T, url string, StatusCode int, Code int, Msg string) (e *httpexpect.Expect) {
	e = getHttpexpect(t)
	e.DELETE(url).WithHeader("Authorization", "Bearer "+GetOauthToken(e)).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Code, Msg)
	return
}

func CreateArticleIp() (*models.ArticleIp, error) {
	m := mock.ArticleIp{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	articleIp := &models.ArticleIp{
		Mun:  m.Mun,
		Addr: m.Addr,
	}
	err = articleIp.CreateArticleIp()
	if err != nil {
		return articleIp, err
	}
	return articleIp, nil
}

func CreateArticle(status string) (*models.Article, error) {
	var tr *models.Type
	var articleIp *models.ArticleIp
	mock.CustomGenerator()
	m := mock.Article{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	if len(status) == 0 {
		status = m.Status
	}

	tr, err = CreateType()
	if err != nil {
		return nil, err
	}

	articleIp, err = CreateArticleIp()
	if err != nil {
		return nil, err
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
		Ips:          []*models.ArticleIp{articleIp},
		Type:         tr,
		TypeID:       tr.ID,
	}
	err = article.CreateArticle()
	if err != nil {
		return article, err
	}
	ArticleCount++
	return article, nil

}

func CreateChapterIp() (*models.ChapterIp, error) {
	m := mock.ChapterIp{}
	err := faker.FakeData(&m)
	if err != nil {
		return nil, err
	}
	chapterIp := &models.ChapterIp{
		Mun:  m.Mun,
		Addr: m.Addr,
	}
	err = chapterIp.CreateChapterIp()
	if err != nil {
		return chapterIp, err
	}
	return chapterIp, nil
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

	var doc *models.Doc
	doc, err = CreateDoc()
	if err != nil {
		return nil, err
	}
	var chapterIp *models.ChapterIp
	chapterIp, err = CreateChapterIp()
	if err != nil {
		return nil, err
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
		Ips:          []*models.ChapterIp{chapterIp},
		Sort:         m.Sort,
		Doc:          doc,
		DocID:        doc.ID,
	}
	err = chapter.CreateChapter()
	if err != nil {
		return chapter, err
	}
	ChapterCount++
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
	ConfigCount++
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
	DocCount++
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
	TypeCount++
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
	TagCount++
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
	PermCount++
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
	RoleCount++
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
	UserCount++
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
