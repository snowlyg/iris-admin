// +build test public type api

package tests

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPublicTypes(t *testing.T) {
	getPublicMore(t, "types", iris.StatusOK, 200, "操作成功")
}
