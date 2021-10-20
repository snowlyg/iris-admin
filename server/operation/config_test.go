package operation

import (
	"reflect"
	"testing"
)

func TestGetExcept(t *testing.T) {
	t.Run("test get except", func(t *testing.T) {
		wantUri := []string{"api/v1/upload", "api/v1/upload"}
		wantMethod := []string{"post", "put"}
		uri, method := CONFIG.GetExcept()
		if !reflect.DeepEqual(uri, wantUri) {
			t.Errorf("get except want %+v but get %+v", wantUri, uri)
		}
		if !reflect.DeepEqual(method, wantMethod) {
			t.Errorf("get except want %+v but get %+v", wantMethod, method)
		}
	})
}
func TestGetInclude(t *testing.T) {
	t.Run("test get include", func(t *testing.T) {
		wantUri := []string{"api/v1/menus"}
		wantMethod := []string{"get"}
		uri, method := CONFIG.GetInclude()
		if !reflect.DeepEqual(uri, wantUri) {
			t.Errorf("get include want %+v but get %+v", wantUri, uri)
		}
		if !reflect.DeepEqual(method, wantMethod) {
			t.Errorf("get include want %+v but get %+v", wantMethod, method)
		}
	})
}
func TestGetIsInclude(t *testing.T) {
	t.Run("test get include", func(t *testing.T) {
		wantUri := "api/v1/menus"
		wantMethod := "get"
		if !CONFIG.IsInclude(wantUri, wantMethod) {
			t.Errorf("[%s](%s) want include true but get false", wantMethod, wantUri)
		}
		wantUri = "api/v1/menus"
		wantMethod = "post"
		if CONFIG.IsInclude(wantUri, wantMethod) {
			t.Errorf("[%s](%s) want include false but get true", wantMethod, wantUri)
		}
		wantUri = "api/v1/menu"
		wantMethod = "get"
		if CONFIG.IsInclude(wantUri, wantMethod) {
			t.Errorf("[%s](%s) want include false but get true", wantMethod, wantUri)
		}
	})
}
func TestGetIsExcept(t *testing.T) {
	t.Run("test get except", func(t *testing.T) {
		wantUri := "api/v1/upload"
		wantMethod := "post"
		if !CONFIG.IsExcept(wantUri, wantMethod) {
			t.Errorf("[%s](%s) want except true but get false", wantMethod, wantUri)
		}
		wantUri = "api/v1/upload"
		wantMethod = "get"
		if CONFIG.IsExcept(wantUri, wantMethod) {
			t.Errorf("[%s](%s) want except false but get true", wantMethod, wantUri)
		}
		wantUri = "api/v1/menu"
		wantMethod = "post"
		if CONFIG.IsExcept(wantUri, wantMethod) {
			t.Errorf("[%s](%s) want except false but get true", wantMethod, wantUri)
		}
	})
}
