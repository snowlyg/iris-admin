// +build test public article api

package tests

import (
	"fmt"
	"github.com/fatih/color"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestPublicArticles(t *testing.T) {
	getPublicMore(t, "article", iris.StatusOK, 200, "操作成功")
}

func TestPublicArticle(t *testing.T) {
	tr, err := CreateArticle("published")
	if err != nil {
		color.Red("TestArticleUpdate %+v", err)
		return
	}

	getOnePublic(t, fmt.Sprintf("article/%d", tr.ID), iris.StatusOK, 200, "操作成功")
}

func TestPublicArticleLike(t *testing.T) {
	tr, err := CreateArticle("published")
	if err != nil {
		color.Red("TestPublicArticleLike %+v", err)
		return
	}

	getOnePublic(t, fmt.Sprintf("article/like/%d", tr.ID), iris.StatusOK, 200, "操作成功")
}
