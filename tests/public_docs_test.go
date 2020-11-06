// +build test public doc api

package tests

import (
	"fmt"
	"github.com/fatih/color"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPublicDocs(t *testing.T) {
	getPublicMore(t, "docs", iris.StatusOK, 200, "操作成功")
}

func TestPublicDoc(t *testing.T) {
	tr, err := CreateDoc()
	if err != nil {
		color.Red("TestDocUpdate %+v", err)
		return
	}

	url := "docs/%d"
	getOnePublic(t, fmt.Sprintf(url, tr.ID), iris.StatusOK, 200, "操作成功")
}
