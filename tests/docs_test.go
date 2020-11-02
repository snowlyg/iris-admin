// +build test doc api

package tests

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestDocs(t *testing.T) {
	getMore(t, "docs", iris.StatusOK, 200, "操作成功")
}

func TestDocCreate(t *testing.T) {
	oj := map[string]interface{}{
		"name": "create_role",
	}

	create(t, "docs", oj, iris.StatusOK, 200, "操作成功")
}

func TestDocUpdate(t *testing.T) {
	tr, err := CreateDoc("tname1")
	if err != nil {
		fmt.Print(err)
		return
	}
	oj := map[string]interface{}{
		"name": "test_update_role",
	}

	url := "docs/%d"
	update(t, fmt.Sprintf(url, tr.ID), oj, iris.StatusOK, 200, "操作成功")
}

func TestDocDelete(t *testing.T) {
	tr, err := CreateDoc("tname2")
	if err != nil {
		fmt.Print(err)
		return
	}
	delete(t, fmt.Sprintf("docs/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
