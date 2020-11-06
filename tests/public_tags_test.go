// +build test public tag api

package tests

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPublicTags(t *testing.T) {
	getPublicMore(t, "tags", iris.StatusOK, 200, "操作成功")
}
