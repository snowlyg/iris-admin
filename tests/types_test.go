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
	getMore(t, "types", iris.StatusOK, 200, "操作成功")
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
