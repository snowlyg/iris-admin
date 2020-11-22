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
	_, err := CreateDoc()
	if err != nil {
		color.Red("TestDocUpdate %+v", err)
		return
	}
	tr2, err := CreateDoc()
	if err != nil {
		color.Red("TestDocUpdate %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": 1, "page": 1, "field": "id,name,created_at"}
	more := &More{tr2.ID, 1, 1, DocCount, []interface{}{"id", "name", "created_at"}}
	getMore(t, "docs", iris.StatusOK, obj, more)
}
func TestDocsLimit(t *testing.T) {
	_, err := CreateDoc()
	if err != nil {
		color.Red("TestDocUpdate %+v", err)
		return
	}
	tr2, err := CreateDoc()
	if err != nil {
		color.Red("TestDocUpdate %+v", nil)
		return
	}
	obj := map[string]interface{}{"limit": 2, "page": 1, "field": "id,name,created_at"}
	more := &More{tr2.ID, 2, 2, DocCount, []interface{}{"id", "name", "created_at"}}
	getMore(t, "docs", iris.StatusOK, obj, more)
}

func TestDocsNoPagination(t *testing.T) {
	_, err := CreateDoc()
	if err != nil {
		color.Red("TestDocUpdate %+v", err)
		return
	}
	tr2, err := CreateDoc()
	if err != nil {
		color.Red("TestDocUpdate %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": -1, "page": -1, "field": "id,name,created_at"}
	more := &More{tr2.ID, -1, DocCount, DocCount, []interface{}{"id", "name", "created_at"}}
	getMore(t, "docs", iris.StatusOK, obj, more)
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
