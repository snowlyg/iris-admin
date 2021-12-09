package test

import (
	"os"
	"testing"

	"github.com/snowlyg/helper/tests"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/rbac/file"
)

var (
	loginUrl  = "/api/v1/auth/login"
	logoutUrl = "/api/v1/users/logout"
	url       = "/api/v1/file"
)

func TestUpload(t *testing.T) {
	if TestServer == nil {
		t.Errorf("TestServer is nil")
	}
	client := TestServer.GetTestLogin(t, loginUrl, nil)
	if client != nil {
		defer client.Logout(logoutUrl, nil)
	} else {
		return
	}
	name := "mysqlPwd.txt"
	fh, err := os.Open("D:/admin/go/src/github.com/snowlyg/iris-admin/web/web_iris/modules/rabc/file/tests/" + name)
	if err != nil {
		t.Error(err)
		return
	}
	defer fh.Close()
	files := []tests.File{
		{
			Key:    "file",
			Path:   name,
			Reader: fh,
		},
	}
	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
		{Key: "data", Value: tests.Responses{
			{Key: "local", Value: file.GetPath(name)},
			{Key: "qiniu", Value: ""},
		}},
	}

	client.UPLOAD(url, pageKeys, files)
}
