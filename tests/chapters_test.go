// +build test chapter api

package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/tests/mock"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestChapters(t *testing.T) {
	tr, err := CreateChapter("")
	if err != nil {
		color.Red("TestChapterUpdate %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": 1, "page": 1, "field": "id,title,created_at"}
	more := &More{tr.ID, 1, 1, ChapterCount, []interface{}{"id", "title", "created_at"}}
	getMore(t, "chapters", iris.StatusOK, obj, more)
}

func TestChaptersNoPagination(t *testing.T) {
	tr, err := CreateChapter("")
	if err != nil {
		color.Red("TestChapterUpdate %+v", err)
		return
	}
	obj := map[string]interface{}{"limit": -1, "page": -1, "field": "id,title,created_at"}
	more := &More{tr.ID, -1, ChapterCount, ChapterCount, []interface{}{"id", "title", "created_at"}}
	getMore(t, "chapters", iris.StatusOK, obj, more)
}

func TestChapterCreate(t *testing.T) {
	mock.CustomGenerator()
	doc, err := CreateDoc()
	if err != nil {
		color.Red("TestChapterCreate %+v", err)
		return
	}
	m := mock.Chapter{}
	err = faker.FakeData(&m)
	if err != nil {
		color.Red("TestChapterCreate %+v", err)
		return
	}
	m.Doc = doc
	m.DocId = int64(doc.ID)

	create(t, "chapters", m, iris.StatusOK, 200, "操作成功")
}

func TestChapterUpdate(t *testing.T) {
	tr, err := CreateChapter("")
	if err != nil {
		color.Red("TestChapterUpdate %+v", err)
		return
	}
	m := mock.Chapter{}
	err = faker.FakeData(&m)
	if err != nil {
		color.Red("TestChapterUpdate %+v", err)
		return
	}
	url := "chapters/%d"
	update(t, fmt.Sprintf(url, tr.ID), m, iris.StatusOK, 200, "操作成功")
}

func TestChapterSetSort(t *testing.T) {
	tr, err := CreateChapter("")
	if err != nil {
		color.Red("TestChapterUpdate %+v", err)
		return
	}
	m := map[string]interface{}{"id": tr.ID, "sort": 100}
	url := "chapters/%d/set_sort"
	update(t, fmt.Sprintf(url, tr.ID), m, iris.StatusOK, 200, "操作成功")
}

func TestChapterChangeSort(t *testing.T) {
	tr, err := CreateChapter("")
	if err != nil {
		color.Red("TestChapterUpdate %+v", err)
		return
	}
	tr2, err := CreateChapter("")
	if err != nil {
		color.Red("TestChapterUpdate %+v", err)
		return
	}

	m := map[string]interface{}{"old_id": tr.ID, "old_sort": tr2.Sort, "new_id": tr2.ID, "new_sort": tr.Sort}
	color.Yellow("%+v", m)
	url := "chapters/sort"
	update(t, url, m, iris.StatusOK, 200, "操作成功")
}

func TestChapterDelete(t *testing.T) {
	tr, err := CreateChapter("")
	if err != nil {
		color.Red("TestChapterDelete %+v", err)
		return
	}
	delete(t, fmt.Sprintf("chapters/%d", tr.ID), iris.StatusOK, 200, "删除成功")
}
