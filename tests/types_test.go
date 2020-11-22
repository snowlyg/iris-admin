// +build test type api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/tests/mock"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestTypes(t *testing.T) {
	tr, err := CreateType()
	if err != nil {
		color.Red("TestTypes %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": 1, "page": 1, "field": "id,name,created_at"}
	more := &More{tr.ID, 1, 1, TypeCount, []interface{}{"name"}}
	getMore(t, "types", iris.StatusOK, obj, more)
}

func TestTypesNoPagination(t *testing.T) {
	tr, err := CreateType()
	if err != nil {
		color.Red("TestTypes %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": -1, "page": -1, "field": "id,name,created_at"}
	more := &More{tr.ID, -1, TypeCount, TypeCount, []interface{}{"id", "name", "created_at"}}
	getMore(t, "types", iris.StatusOK, obj, more)
}

func TestTypeCreate(t *testing.T) {
	m := mock.Type{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestTypeCreate %+v", err)
		return
	}

	create(t, "types", m, iris.StatusOK, 200, "操作成功")
}

func TestTypeUpdate(t *testing.T) {
	tr, err := CreateType()
	if err != nil {
		color.Red("TestTypeUpdate %+v", err)
		return
	}
	m := mock.Type{}
	err = faker.FakeData(&m)
	if err != nil {
		color.Red("TestTypeUpdate %+v", err)
		return
	}

	url := "types/%d"
	update(t, fmt.Sprintf(url, tr.ID), m, iris.StatusOK, 200, "操作成功")
}

func TestTypeDelete(t *testing.T) {
	tr, err := CreateType()
	if err != nil {
		fmt.Print(err)
		return
	}
	delete(t, fmt.Sprintf("types/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
