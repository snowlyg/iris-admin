// +build test tag api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/tests/mock"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestTags(t *testing.T) {
	getMore(t, "tags", iris.StatusOK, 200, "操作成功")
}

func TestTagCreate(t *testing.T) {
	m := mock.Tag{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestTagCreate %+v", err)
		return
	}

	create(t, "tags", m, iris.StatusOK, 200, "操作成功")
}

func TestTagUpdate(t *testing.T) {
	m := mock.Tag{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestTagUpdate %+v", err)
		return
	}

	tr, err := CreateTag()
	if err != nil {
		color.Red("TestTagUpdate %+v", err)
		return
	}

	url := "tags/%d"
	update(t, fmt.Sprintf(url, tr.ID), m, iris.StatusOK, 200, "操作成功")
}

func TestTagDelete(t *testing.T) {
	tr, err := CreateTag()
	if err != nil {
		fmt.Print(err)
		return
	}
	delete(t, fmt.Sprintf("tags/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
