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
	tr, err := CreateTag()
	if err != nil {
		color.Red("TestTags %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": 1, "page": 1, "field": "id,name,created_at"}
	more := &More{tr.ID, 1, 1, 1, []interface{}{"id", "name", "created_at"}}
	getMore(t, "tags", iris.StatusOK, obj, more)
}

func TestTagsNoPagination(t *testing.T) {
	tr, err := CreateTag()
	if err != nil {
		color.Red("TestTagsNoPagination %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": -1, "page": -1, "field": "id,name,created_at"}
	more := &More{tr.ID, -1, TagCount, TagCount, []interface{}{"id", "name", "created_at"}}
	getMore(t, "tags", iris.StatusOK, obj, more)
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
