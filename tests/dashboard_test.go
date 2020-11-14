// +build test dashboard api

package tests

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestDashBoard(t *testing.T) {
	getData(t, "dashboard", iris.StatusOK, nil, []interface{}{
		"articleVisitis", "articles", "docVisitis", "docs",
	})
}
