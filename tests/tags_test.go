// +build test tag api

package tests

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestTags(t *testing.T) {
	getMore(t, "tags", iris.StatusOK, 200, "操作成功")
}

func TestTagCreate(t *testing.T) {
	oj := map[string]interface{}{
		"name": "create_role",
	}

	create(t, "tags", oj, iris.StatusOK, 200, "操作成功")
}

func TestTagUpdate(t *testing.T) {
	tr, err := CreateTag("tname1")
	if err != nil {
		fmt.Print(err)
		return
	}
	oj := map[string]interface{}{
		"name": "test_update_role",
	}

	url := "tags/%d"
	update(t, fmt.Sprintf(url, tr.ID), oj, iris.StatusOK, 200, "操作成功")
}

func TestTagDelete(t *testing.T) {
	tr, err := CreateTag("tname2")
	if err != nil {
		fmt.Print(err)
		return
	}
	delete(t, fmt.Sprintf("tags/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
