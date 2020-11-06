// +build test public chapter api

package tests

import (
	"fmt"
	"github.com/fatih/color"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPublicChapters(t *testing.T) {
	getPublicMore(t, "chapter", iris.StatusOK, 200, "操作成功")
}

func TestPublicChapter(t *testing.T) {
	tr, err := CreateChapter("published")
	if err != nil {
		color.Red("TestChapterUpdate %+v", err)
		return
	}

	getOnePublic(t, fmt.Sprintf("chapter/%d", tr.ID), iris.StatusOK, 200, "操作成功")
}

func TestPublicChapterLike(t *testing.T) {
	tr, err := CreateChapter("published")
	if err != nil {
		color.Red("TestPublicChapterLike %+v", err)
		return
	}

	getOnePublic(t, fmt.Sprintf("chapter/like/%d", tr.ID), iris.StatusOK, 200, "操作成功")
}
