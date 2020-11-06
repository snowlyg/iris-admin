// +build test public config api

package tests

import (
	"fmt"
	"github.com/fatih/color"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPublicConfig(t *testing.T) {
	tr, err := CreateConfig()
	if err != nil {
		color.Red("TestPublicConfig %+v", err)
		return
	}

	getOnePublic(t, fmt.Sprintf("config/%s", tr.Name), iris.StatusOK, 200, "操作成功")
}
