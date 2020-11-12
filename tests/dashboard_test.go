// +build test dashboard api

package tests

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestDashBoard(t *testing.T) {
	getMore(t, "dashboard", iris.StatusOK, 200, "操作成功")
}
