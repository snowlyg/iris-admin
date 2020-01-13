package main

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPermissions(t *testing.T) {
	getMore(t, baseUrl+"permissions", iris.StatusOK, true, "操作成功")
}

func TestImportPermissions(t *testing.T) {
	bImport(t, baseUrl+"permissions/import", iris.StatusOK, true, "操作成功", nil)
}
