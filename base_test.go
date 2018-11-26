package main

import (
	"IrisYouQiKangApi/logic"
	"IrisYouQiKangApi/system"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
	"github.com/kataras/iris/sessions"
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
	system.Redis.Set("iris", sessions.LifeTime{}, "project_env", "testing", false)
	App = NewApp()
	m.Run()
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
		system.Config.Get("test.LoginUser").(string),
		system.Config.Get("test.LoginPwd").(string),
	)

	if bc.Data != nil {
		e.GET(bc.Url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(bc.StatusCode).JSON().Object().Values().Contains(bc.Status, bc.Msg, bc.Data)
	} else {
		e.GET(bc.Url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(bc.StatusCode).JSON().Object().Values().
			Contains(bc.Status, bc.Msg).
			Contains().Last().Object().Keys().Contains(
			"created_at",
			"deleted_at",
			"email",
			"id",
			"is_audit",
			"is_client",
			"is_client_admin",
			"is_frozen",
			"is_wechat_verfiy",
			"name",
			"open_id",
			"phone",
			"remember_token",
			"role_id",
			"role_name",
			"updated_at",
			"username",
			"wechat_avatar",
			"wechat_name",
			"wechat_verfiy_time",
		)
		//Element(2).Array().Length().Equal(20)
	}

	return
}
