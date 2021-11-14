package test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/snowlyg/iris-admin/server/web/web_gin"
	v1 "github.com/snowlyg/iris-admin/server/web/web_gin/modules/v1"
)

//go:embed mysqlPwd.txt
var mysqlPwd string

//go:embed redisPwd.txt
var redisPwd string

var TestServer *web_gin.WebServer

func TestMain(m *testing.M) {
	var uuid string
	uuid, TestServer = v1.BeforeTestMain(mysqlPwd, redisPwd, 2)
	code := m.Run()
	v1.AfterTestMain(uuid)

	os.Exit(code)
}
