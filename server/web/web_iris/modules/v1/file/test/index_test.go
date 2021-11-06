package test

import (
	"testing"

	"github.com/snowlyg/helper/tests"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/file"
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
	defer client.Logout(logoutUrl, nil)
	files := map[string]string{
		"file": "./avatar.jpg",
	}
	filename, err := file.GetFileName(files["file"])
	if err != nil {
		t.Fatalf(err.Error())
	}
	pageKeys := tests.Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
		{Key: "data", Value: tests.Responses{
			{Key: "local", Value: file.GetPath(filename)},
			{Key: "qiniu", Value: ""},
		}},
	}

	client.UPLOAD(url, pageKeys, files)
}
