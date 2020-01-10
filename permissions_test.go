package main

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPermissions(t *testing.T) {
	url := "/v1/admin/permissions"
	getMore(t, url, iris.StatusOK, true, "操作成功", nil)
}
