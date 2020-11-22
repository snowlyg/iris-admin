// +build test config api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/tests/mock"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestConfigs(t *testing.T) {
	getAll(t, "configs", iris.StatusOK, nil, ConfigCount)
}

func TestConfigCreate(t *testing.T) {
	m := mock.Config{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestConfigCreate %+v", err)
		return
	}

	create(t, "configs", m, iris.StatusOK, 200, "操作成功")
}

func TestConfigUpdate(t *testing.T) {
	m := mock.Config{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestConfigUpdate %+v", err)
		return
	}
	tr, err := CreateConfig()
	if err != nil {
		color.Red("TestConfigUpdate %+v", err)
		return
	}

	url := "configs/%d"
	update(t, fmt.Sprintf(url, tr.ID), m, iris.StatusOK, 200, "操作成功")
}

func TestConfigGetByKey(t *testing.T) {
	tr, err := CreateConfig()
	if err != nil {
		fmt.Print(err)
		return
	}
	getOne(t, fmt.Sprintf("configs/%s", tr.Name), iris.StatusOK, 200, "操作成功")
}

func TestConfigDelete(t *testing.T) {
	tr, err := CreateConfig()
	if err != nil {
		fmt.Print(err)
		return
	}
	delete(t, fmt.Sprintf("configs/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
