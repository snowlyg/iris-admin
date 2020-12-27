// +build test

package libs

import (
	"reflect"
	"testing"
)

func TestStructToMap(t *testing.T) {
	type Js struct {
		Name  string
		Title string
	}

	js := Js{"name", "title"}
	want := map[string]interface{}{
		"name":  "name",
		"title": "title",
	}

	t.Run("TestStructToMap", func(t *testing.T) {
		if got := StructToMap(js); !reflect.DeepEqual(got, want) {
			t.Errorf("StructToMap() = %v, want %v", got, want)
		}
	})
}

func TestStructToString(t *testing.T) {
	type Js struct {
		Name  string
		Title string
	}

	js := Js{"name", "title"}
	want := `{"Name":"name","Title":"title"}`

	t.Run("TestStructToString", func(t *testing.T) {
		if got := StructToString(js); got != want {
			t.Errorf("StructToString() = %v, want %v", got, want)
		}
	})

}
