// +build test

package main

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPermissions(t *testing.T) {
	getMore(t, "permissions", iris.StatusOK, 200, "操作成功")
}
