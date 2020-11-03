// +build test expire api

package tests

import (
	"testing"

	"github.com/kataras/iris/v12"
)

func TestExpire(t *testing.T) {
	getMore(t, "expire", iris.StatusOK, 200, "")
}
