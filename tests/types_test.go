// +build test type api

package tests

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestTypes(t *testing.T) {
	getMore(t, "types", iris.StatusOK, 200, "操作成功")
}

func TestTypeCreate(t *testing.T) {
	oj := map[string]interface{}{
		"name": "create_role",
	}

	create(t, "types", oj, iris.StatusOK, 200, "操作成功")
}

func TestTypeUpdate(t *testing.T) {
	tr, err := CreateType("tname1")
	if err != nil {
		fmt.Print(err)
		return
	}
	oj := map[string]interface{}{
		"name": "test_update_role",
	}

	url := "types/%d"
	update(t, fmt.Sprintf(url, tr.ID), oj, iris.StatusOK, 200, "操作成功")
}

func TestTypeDelete(t *testing.T) {
	tr, err := CreateType("tname2")
	if err != nil {
		fmt.Print(err)
		return
	}
	delete(t, fmt.Sprintf("types/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
