package main

import (
	"IrisYouQiKangApi/config"
	"IrisYouQiKangApi/database"
	Redis "IrisYouQiKangApi/redis"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
	"testing"
)

type BaseCase struct {
	Url        string      //测试路由
	Object     interface{} //发送的json 对象
	StatusCode int         //返回的 http 状态码
	Status     bool        //返回的状态
	Msg        string      //返回提示信息
	Data       interface{} //返回数据
}

var (
	app           *iris.Application //iris.Applications
	rd            *redis.Client     //iris.Applications
	testAdminUser *Users
)

func TestMain(m *testing.M) {
	//初始化配置
	conf = config.New()
	//初始化redis
	rd = Redis.New()
	//初始化测试数据库
	db = database.New(conf, "testing")

	//获取测试的数据表
	value, err := rd.Get("test_table_name").Result()
	if err == redis.Nil {
		fmt.Println("env_t does not exist")
	} else if err != nil {
		panic(err)
	}

	app = NewApp()

	testAdminUser = CreaterSystemAdmin()

	m.Run()

	//删除测试数据表，保持测试环境
	db.DropTable(value)
}

//单元测试 post 方法
func (bc *BaseCase) login(t *testing.T) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: false})
	if bc.Data != nil {
		e.POST(bc.Url).WithJSON(bc.Object).
			Expect().Status(bc.StatusCode).JSON().Object().Values().Contains(bc.Status, bc.Msg, bc.Data)
	} else {
		e.POST(bc.Url).WithJSON(bc.Object).
			Expect().Status(bc.StatusCode).JSON().Object().Values().Contains(bc.Status, bc.Msg)
	}

	return
}

//单元测试 post 方法
func (bc *BaseCase) post(t *testing.T) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: true})
	at := GetLoginToken()

	fmt.Println(at)
	if bc.Data != nil {
		e.POST(bc.Url).WithHeader("Authorization", "Bearer "+at.Token).WithJSON(bc.Object).
			Expect().Status(bc.StatusCode).JSON().Object().Values().Contains(bc.Status, bc.Msg, bc.Data)
	} else {
		e.POST(bc.Url).WithHeader("Authorization", "Bearer "+at.Token).WithJSON(bc.Object).
			Expect().Status(bc.StatusCode).JSON().Object().Values().Contains(bc.Status, bc.Msg)
	}

	return
}

//单元测试 get 方法
func (bc *BaseCase) get(t *testing.T) (e *httpexpect.Expect) {
	e = httptest.New(t, app)
	at := GetLoginToken()
	if bc.Data != nil {
		e.GET(bc.Url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(bc.StatusCode).JSON().Object().Values().Contains(bc.Status, bc.Msg, bc.Data)
	} else {
		e.GET(bc.Url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(bc.StatusCode).JSON().Object().Values().
			Contains(bc.Status, bc.Msg)
	}

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
	aul := new(AdminUserLogin)
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

	//打印错误信息
	if !status {
		fmt.Println(msg)
	}

	return response
}
