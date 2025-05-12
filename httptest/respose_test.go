package httptest

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/helper/arr"
)

func TestIdKeys(t *testing.T) {
	want := Responses{
		{Key: "id", Value: 0, Type: "ge"},
	}
	t.Run("Test id keys", func(t *testing.T) {
		idKeys := IdKeys()
		if !reflect.DeepEqual(want, idKeys) {
			t.Errorf("IdKeys want %+v but get %+v", want, idKeys)
		}
	})
}

func TestHttpTest(t *testing.T) {
	engine := gin.New()
	// Add /example route via handler function to the gin instance
	handler := GinHandler(engine)
	// Create httpexpect instance
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	pageKeys := Responses{{Key: "message", Value: "OK"}, {Key: "status", Value: 200}, {Key: "data", Value: Response{Key: "message", Value: "pong"}}}
	value := e.GET("/example").Expect().Status(http.StatusOK).JSON()

	Test(value, pageKeys)
}

func TestHttpTestArray(t *testing.T) {
	engine := gin.New()
	// Add /example route via handler function to the gin instance
	handler := GinHandler(engine)
	// Create httpexpect instance
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	pageKeys := []string{"1", "2"}
	value := e.GET("/array").Expect().Status(http.StatusOK).JSON()
	Test(value, pageKeys)
}

func TestHttpScan(t *testing.T) {
	engine := gin.New()
	// Add /example route via handler function to the gin instance
	handler := GinHandler(engine)
	// Create httpexpect instance
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	pageKeys := Responses{{Key: "message", Value: ""}}
	obj := e.GET("/example").Expect().Status(http.StatusOK).JSON().Object()

	Scan(obj, pageKeys)
	x := pageKeys.GetString("data.message")
	if x != "pong" {
		t.Errorf("Scan want get pong but get %s", x)
	}
}

func TestSchema(t *testing.T) {
	wantJson := `{
		"status": 200,
		"data": {
			"list": [
				{
					"createdAt": "2025-03-21T16:27:20+08:00",
					"deletedAt": "",
					"updatedAt": "2025-03-21T16:27:20+08:00",
					"dev_remark": "",
					"pac_room_id": 1,
					"room_desc": "1413-301"
				}
			],
			"total": 1,
			"page": 1,
			"pageSize": 20
		},
		"message": "OK"
	}`
	res, err := Schema([]byte(wantJson))
	if err != nil {
		t.Fatal(err.Error())
	}
	if res.GetInt("status") != 200 {
		t.Errorf("status want %d but get %d", 200, res.GetInt("status"))
	}
	if res.GetString("message") != "OK" {
		t.Errorf("message want %s but get %s", "OK", res.GetString("message"))
	}
	data := res.GetResponse("data")
	if data != nil {
		if data.GetInt("total") != 1 {
			t.Errorf("total want %d but get %d", 1, data.GetInt("total"))
		}
		if data.GetInt("page") != 1 {
			t.Errorf("page want %d but get %d", 1, data.GetInt("page"))
		}
		if data.GetInt("pageSize") != 20 {
			t.Errorf("pageSize want %d but get %d", 20, data.GetInt("pageSize"))
		}
		list := data.GetResponses("list")
		if len(list) != 1 {
			t.Errorf("list len want %d but get %d", 1, len(list))
		}
		first := list[0]
		if first.GetId("pac_room_id") != 1 {
			t.Errorf("pac_room_id want %d but get %d", 1, first.GetId("pac_room_id"))
		}
		roomDesc := first.GetString("room_desc")
		if roomDesc != "1413-301" {
			t.Errorf("room_desc want %s but get '%s'", "1413-301", roomDesc)
		}
		if first.GetString("dev_remark") != "" {
			t.Errorf("dev_remark want %s but get '%s'", "", first.GetString("dev_remark"))
		}
		keys := arr.NewCheckArrayType(0)
		for _, v := range first.Keys() {
			keys.Add(v)
		}

		for _, k := range []string{"createdAt", "deletedAt", "updatedAt", "dev_remark", "pac_room_id", "room_desc"} {
			if !keys.Check(k) {
				t.Errorf("%s not in keys", k)
			}
		}
	}
}

func TestSchemaResponse(t *testing.T) {
	data := `{
		"status": 200,
		"message": "OK"
	}`
	j := map[string]any{}
	if err := json.Unmarshal([]byte(data), &j); err != nil {
		t.Error(err.Error())
	}
	log.Printf("j %+v\n", j)
	wantKey := "data"
	resp := schemaResponse(wantKey, j)

	if value, ok := resp.Value.(Responses); !ok {
		t.Error("schema response return value not Responses")
	} else {
		keys := arr.NewCheckArrayType(0)
		for _, v := range value {
			keys.Add(v.Key)
			if v.Key == "message" {
				wantValue := "OK"
				if v.Value != wantValue {
					t.Errorf("%s Value want '%v' but get '%v'", v.Key, wantValue, v.Value)
				}
			} else if v.Key == "status" {
				var wantValue float64 = 200
				if v.Value != wantValue {
					t.Errorf("%s Value want '%v' but get '%v'", v.Key, wantValue, v.Value)
				}
			} else {
				t.Errorf("key %s is in response", v.Key)
			}
		}
		if !keys.Check("message") {
			t.Error("message not in keys")
		}
		if !keys.Check("status") {
			t.Error("status not in keys")
		}
	}
}
