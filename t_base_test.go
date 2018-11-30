package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"IrisApiProject/config"
	"IrisApiProject/database"
	Redis "IrisApiProject/redis"
	"github.com/go-redis/redis"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
)

var (
	app           *iris.Application // iris.Applications
	rd            *redis.Client     // iris.Applications
	testAdminUser *Users
)

func TestMain(m *testing.M) {
	// 初始化配置
	conf = config.New()
	// 初始化redis
	rd = Redis.New()
	// 初始化测试数据库
	db = database.New(conf, "testing")

	// 获取测试的数据表
	value, err := rd.Get("test_table_name").Result()
	if err == redis.Nil {
		fmt.Println("env_t does not exist")
	} else if err != nil {
		panic(err)
	}

	// 初始化app
	app = NewApp()
	// 创建用户
	testAdminUser = CreaterSystemAdmin()

	flag.Parse()
	exitCode := m.Run()

	// 删除测试数据表，保持测试环境
	db.DropTable(value)

	os.Exit(exitCode)
}

// 单元测试 post 方法
func login(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: conf.Get("app.debug").(bool)})
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

// 单元测试 post 方法
func create(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: conf.Get("app.debug").(bool)})
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

// 单元测试 post 方法
func update(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: conf.Get("app.debug").(bool)})
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

// 单元测试 get 方法
func getOne(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: conf.Get("app.debug").(bool)})
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

// 单元测试 get 方法
func getMore(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: conf.Get("app.debug").(bool)})
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

// 单元测试 get 方法
func delete(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: conf.Get("app.debug").(bool)})
	at := GetLoginToken()

	e.DELETE(url).WithHeader("Authorization", "Bearer "+at.Token).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)

	return
}

/**
*设置测试数据表
*@param tn stirng 数据表名称
 */
func SetTestTableName(tn string) {
	err := rd.Set("test_table_name", tn, 0).Err()
	if err != nil {
		panic(err)
	}
}

/**
*创建系统管理员
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreaterSystemAdmin() *Users {
	aul := new(UserJson)
	aul.Username = conf.Get("test.LoginUserName").(string)
	aul.Password = conf.Get("test.LoginPwd").(string)
	aul.Phone = "12345678"
	aul.Name = conf.Get("test.LoginName").(string)
	aul.RoleId = 1

	return MCreateUser(aul)
}

/**
*登陆用户
*@return   Token 返回登陆后的token
 */
func GetLoginToken() Token {
	response, status, msg := LUserAdminCheckLogin(
		conf.Get("test.LoginUserName").(string),
		conf.Get("test.LoginPwd").(string),
	)

	// 打印错误信息
	if !status {
		fmt.Println(msg)
	}

	return response
}
