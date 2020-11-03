// +build test doc api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/tests/mock"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestDocs(t *testing.T) {
	getMore(t, "docs", iris.StatusOK, 200, "操作成功")
}

func TestDocCreate(t *testing.T) {
	m := mock.Doc{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestDocCreate %+v", err)
		return
	}

	create(t, "docs", m, iris.StatusOK, 200, "操作成功")
}

func TestDocUpdate(t *testing.T) {
	m := mock.Doc{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestDocUpdate %+v", err)
		return
	}
	tr, err := CreateDoc()
	if err != nil {
		color.Red("TestDocUpdate %+v", err)
		return
	}

	url := "docs/%d"
	update(t, fmt.Sprintf(url, tr.ID), m, iris.StatusOK, 200, "操作成功")
}

func TestDocDelete(t *testing.T) {
	tr, err := CreateDoc()
	if err != nil {
		fmt.Print(err)
		return
	}
	delete(t, fmt.Sprintf("docs/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
