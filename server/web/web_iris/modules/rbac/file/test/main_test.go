package test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/snowlyg/helper/tests"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
	v1 "github.com/snowlyg/iris-admin/server/web/web_iris/modules/rbac"
)

//go:embed mysqlPwd.txt
var mysqlPwd string

//go:embed redisPwd.txt
var redisPwd string

var TestServer *web_iris.WebServer
var TestClient *tests.Client

func TestMain(m *testing.M) {
	var uuid string
	uuid, TestServer = v1.BeforeTestMain(mysqlPwd, redisPwd, 1)
	code := m.Run()
	v1.AfterTestMain(uuid, TestClient)

	os.Exit(code)
}
