package main

import (
	"IrisYouQiKangApi/logic"
	"IrisYouQiKangApi/models"
	"IrisYouQiKangApi/system"
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
	App *iris.Application //iris.Applications
)

func TestMain(m *testing.M) {
	//删除测试数据表，保持测试环境
	ttn := system.RedisGet("test_table_name")
	App = NewApp()

	//不是 users 表测试，自动创建系统管理员
	if ttn != "users" {
		CreaterSystemAdmin()
	}

	m.Run()

	system.DB.DropTable(ttn)
}

//单元测试 post 方法
func (bc *BaseCase) post(t *testing.T) (e *httpexpect.Expect) {
	e = httptest.New(t, App, httptest.Configuration{Debug: false})
	if bc.Data != nil {
		e.POST(bc.Url).WithJSON(bc.Object).
			Expect().Status(bc.StatusCode).JSON().Object().Values().Contains(bc.Status, bc.Msg, bc.Data)
	} else {
		e.POST(bc.Url).WithJSON(bc.Object).
			Expect().Status(bc.StatusCode).JSON().Object().Values().Contains(bc.Status, bc.Msg)
	}

	return
}

//单元测试 get 方法
func (bc *BaseCase) get(t *testing.T) (e *httpexpect.Expect) {
	e = httptest.New(t, App)
	at, _, _ := logic.UserAdminCheckLogin(
		system.Config.Get("test.LoginUserName").(string),
		system.Config.Get("test.LoginPwd").(string),
	)

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
	system.RedisSet("test_table_name", tn, 0)
}

/**
*创建系统管理员
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreaterSystemAdmin() *models.Users {
	aul := new(models.AdminUserLogin)
	aul.Username = system.Config.Get("test.LoginUserName").(string)
	aul.Password = system.Config.Get("test.LoginPwd").(string)
	aul.Phone = "12345678"
	aul.Name = system.Config.Get("test.LoginName").(string)
	aul.RoleId = 1

	return models.CreateUser(aul)
}
