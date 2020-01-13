package main

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPermissions(t *testing.T) {
	getMore(t, "/v1/admin/permissions", iris.StatusOK, true, "操作成功", nil)
}

//func TestImportPermissions(t *testing.T) {
//	bImport(t, "/v1/admin/permissions/import", iris.StatusOK, true, "操作成功", nil)
//}
