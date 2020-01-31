package main

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestResetData(t *testing.T) {
	getOnAuth(t, baseUrl+"resetData", iris.StatusOK, true, "重置数据成功")
}
