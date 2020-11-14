// +build test article api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/tests/mock"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestArticles(t *testing.T) {
	_, err := CreateArticle("")
	if err != nil {
		color.Red("TestArticles %+v", err)
		return
	}
	tr2, err := CreateArticle("")
	if err != nil {
		color.Red("TestArticles %+v", err)
		return
	}

	obj := map[string]interface{}{"limit": 1, "page": 1}
	more := &More{tr2.ID, 1, 1, 2}
	getMore(t, "article", iris.StatusOK, obj, more)
}

func TestArticlesWithSortByAsc(t *testing.T) {
	_, err := CreateArticle("")
	if err != nil {
		color.Red("TestArticles %+v", err)
		return
	}
	_, err = CreateArticle("")
	if err != nil {
		color.Red("TestArticles %+v", err)
		return
	}

	obj := map[string]interface{}{"limit": 1, "page": 1, "sort": "asc"}
	more := &More{1, 1, 1, 4}
	getMore(t, "article", iris.StatusOK, obj, more)
}

func TestArticlesWithNoPagination(t *testing.T) {
	_, err := CreateArticle("")
	if err != nil {
		color.Red("TestArticles %+v", err)
		return
	}
	_, err = CreateArticle("")
	if err != nil {
		color.Red("TestArticles %+v", err)
		return
	}

	obj := map[string]interface{}{"limit": -1, "page": -1, "sort": "asc"}
	more := &More{1, -1, 6, 6}
	getMore(t, "article", iris.StatusOK, obj, more)
}

func TestArticleCreate(t *testing.T) {
	mock.CustomGenerator()
	m := mock.Article{}
	err := faker.FakeData(&m)
	if err != nil {
		color.Red("TestArticleCreate %+v", err)
		return
	}

	create(t, "article", m, iris.StatusOK, 200, "操作成功")
}

func TestArticleUpdate(t *testing.T) {
	tr, err := CreateArticle("")
	if err != nil {
		color.Red("TestArticleUpdate %+v", err)
		return
	}
	m := mock.Article{}
	err = faker.FakeData(&m)
	if err != nil {
		color.Red("TestArticleUpdate %+v", err)
		return
	}
	url := "article/%d"
	update(t, fmt.Sprintf(url, tr.ID), m, iris.StatusOK, 200, "操作成功")
}

func TestArticleDelete(t *testing.T) {
	tr, err := CreateArticle("")
	if err != nil {
		color.Red("TestArticleDelete %+v", err)
		return
	}
	delete(t, fmt.Sprintf("article/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
