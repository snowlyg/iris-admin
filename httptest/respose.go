package httptest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gavv/httpexpect/v2"
	"github.com/snowlyg/helper/arr"
)

// Responses
type Responses []Response

// Response
type Response struct {
	Type   string        // httpest type , if empty use  IsEqual() function to test
	Key    string        // httptest data's key
	Value  any           // httptest data's value
	Length int           // httptest data's length,when the data are array or map
	Func   func(obj any) // httpest func, you can add your test logic ,can be empty
}

// Keys return Responses object key array
func (res Responses) Keys() []string {
	keys := []string{}
	for _, re := range res {
		keys = append(keys, re.Key)
	}
	return keys
}

// IdKeys return Responses with id
func IdKeys() Responses {
	return Responses{
		{Key: "id", Value: 0, Type: "ge"},
	}
}

// Test for data test
func Test(value *httpexpect.Value, reses ...any) {
	for _, ks := range reses {
		if ks == nil {
			return
		}
		reflectTypeString := reflect.TypeOf(ks).String()
		switch reflectTypeString {
		case "bool":
			value.Boolean().IsEqual(ks.(bool))
		case "string":
			value.String().IsEqual(ks.(string))
		case "float64":
			value.Number().IsEqual(ks.(float64))
		case "uint":
			value.Number().IsEqual(ks.(uint))
		case "int":
			value.IsEqual(ks.(int))

		case "[]httptest.Responses":
			valueLen := len(ks.([]Responses))
			length := int(value.Array().Length().Raw())
			value.Array().Length().IsEqual(valueLen)
			if length > 0 {
				max := 1
				if valueLen == length {
					max = length
				}
				for i := 0; i < max; i++ {
					ks.([]Responses)[i].Test(value.Array().Value(i))
				}
			}

		case "map[int][]httptest.Responses":
			values := ks.(map[int][]Responses)
			length := len(values)
			value.Object().Keys().Length().IsEqual(length)
			if length > 0 {
				for key, v := range values {
					for _, vres := range v {
						vres.Test(value.Object().Value(strconv.FormatInt(int64(key), 10)))
					}
				}
			}
		case "httptest.Responses":
			ks.(Responses).Test(value)
		case "[]uint":
			valueLen := len(ks.([]uint))
			value.Array().Length().IsEqual(valueLen)
			length := int(value.Array().Length().Raw())
			if length > 0 {
				max := 1
				if valueLen == length {
					max = length
				}
				for i := 0; i < max; i++ {
					value.Array().Value(i).Number().IsEqual(ks.([]uint)[i])
				}
			}

		case "[]string":
			valueLen := len(ks.([]string))
			value.Array().Length().IsEqual(valueLen)
			length := int(value.Array().Length().Raw())
			if length > 0 {
				max := 1
				if valueLen == length {
					max = length
				}
				for i := 0; i < max; i++ {
					value.Array().Value(i).String().IsEqual(ks.([]string)[i])
				}
			}
		case "map[int]string":
			values := ks.(map[int]string)
			value.Object().Keys().Length().IsEqual(len(values))
			for key, v := range values {
				value.Object().Value(strconv.FormatInt(int64(key), 10)).IsEqual(v)
			}
		default:
			continue
		}
	}
}

// Scan scan data form http response
func Scan(object *httpexpect.Object, reses ...Responses) {
	if len(reses) == 0 {
		return
	}

	//return once
	if len(reses) == 1 {
		reses[0].Scan(object.Value("data").Object())
		return
	}

	array := object.Value("data").Array()
	length := int(array.Length().Raw())
	if length < len(reses) {
		fmt.Println("Return data not IsEqual keys length")
		array.Length().IsEqual(len(reses))
		return
	}

	// return array
	for m, res := range reses {
		if res == nil {
			return
		}
		res.Scan(object.Value("data").Array().Value(m).Object())
	}
}

// Test Test Responses object
func (resp Responses) Test(value *httpexpect.Value) {
	for _, rs := range resp {
		if rs.Value == nil {
			continue
		}
		if rs.Func != nil {
			rs.Func(value.Object().Value(rs.Key))

		} else {
			reflectTypeString := reflect.TypeOf(rs.Value).String()
			switch reflectTypeString {
			case "bool":
				value.Object().Value(rs.Key).Boolean().IsEqual(rs.Value.(bool))
			case "string":
				if strings.ToLower(rs.Type) == "notempty" {
					value.Object().Value(rs.Key).String().NotEmpty()
				} else {
					value.Object().Value(rs.Key).String().IsEqual(rs.Value.(string))
				}
			case "float64":
				if strings.ToLower(rs.Type) == "ge" {
					value.Object().Value(rs.Key).Number().Ge(rs.Value.(float64))
				} else {
					value.Object().Value(rs.Key).Number().IsEqual(rs.Value.(float64))
				}
			case "uint":
				if strings.ToLower(rs.Type) == "ge" {
					value.Object().Value(rs.Key).Number().Ge(rs.Value.(uint))
				} else {
					value.Object().Value(rs.Key).Number().IsEqual(rs.Value.(uint))
				}
			case "int":
				if strings.ToLower(rs.Type) == "ge" {
					value.Object().Value(rs.Key).Number().Ge(rs.Value.(int))
				} else {
					value.Object().Value(rs.Key).Number().IsEqual(rs.Value.(int))
				}
			case "[]httptest.Responses":
				valueLen := len(rs.Value.([]Responses))
				length := int(value.Object().Value(rs.Key).Array().Length().Raw())
				value.Object().Value(rs.Key).Array().Length().IsEqual(valueLen)
				if length > 0 {
					max := 1
					if rs.Length > 0 {
						max = rs.Length
					}
					if valueLen == length {
						max = length
					}
					if valueLen > 0 {
						for i := 0; i < max; i++ {
							rs.Value.([]Responses)[i].Test(value.Object().Value(rs.Key).Array().Value(i))
						}
					}
				}

			case "map[int][]httptest.Responses":
				values := rs.Value.(map[int][]Responses)
				length := len(values)
				value.Object().Value(rs.Key).Object().Keys().Length().IsEqual(length)
				if length > 0 {
					for key, v := range values {
						for _, vres := range v {
							vres.Test(value.Object().Value(rs.Key).Object().Value(strconv.FormatInt(int64(key), 10)))
						}
					}
				}
			case "httptest.Responses":
				rs.Value.(Responses).Test(value.Object().Value(rs.Key))
			case "[]uint":
				valueLen := len(rs.Value.([]uint))
				value.Object().Value(rs.Key).Array().Length().IsEqual(valueLen)
				length := int(value.Object().Value(rs.Key).Array().Length().Raw())
				if length > 0 {
					max := 1
					if rs.Length > 0 {
						max = rs.Length
					}
					if valueLen == length {
						max = length
					}
					for i := 0; i < max; i++ {
						value.Object().Value(rs.Key).Array().ContainsAny(rs.Value.([]uint)[i])
					}
				}

			case "[]string":

				if strings.ToLower(rs.Type) == "null" {
					value.Object().Value(rs.Key).IsNull()
				} else if strings.ToLower(rs.Type) == "notnull" {
					value.Object().Value(rs.Key).NotNull()
				} else {
					valueLen := len(rs.Value.([]string))
					value.Object().Value(rs.Key).Array().Length().IsEqual(valueLen)
					length := int(value.Object().Value(rs.Key).Array().Length().Raw())
					if length > 0 {
						max := 1
						if rs.Length > 0 {
							max = rs.Length
						}
						if valueLen == length {
							max = length
						}
						for i := 0; i < max; i++ {
							value.Object().Value(rs.Key).Array().ContainsAny(rs.Value.([]string)[i])
						}
					}
				}
			case "map[int]string":
				if strings.ToLower(rs.Type) == "null" {
					value.Object().Value(rs.Key).IsNull()
				} else if strings.ToLower(rs.Type) == "notnull" {
					value.Object().Value(rs.Key).NotNull()
				} else {
					values := rs.Value.(map[int]string)
					value.Object().Value(rs.Key).Object().Keys().Length().IsEqual(len(values))
					for key, v := range values {
						value.Object().Value(rs.Key).Object().Value(strconv.FormatInt(int64(key), 10)).IsEqual(v)
					}
				}
			default:
				continue
			}
		}
	}
	resp.Scan(value.Object())
}

// Scan Scan response data to Responses object.
func (res Responses) Scan(object *httpexpect.Object) {
	for k, rk := range res {
		if !Exist(object, rk.Key) {
			continue
		}
		if rk.Value == nil {
			continue
		}
		valueTypeName := reflect.TypeOf(rk.Value).String()
		switch valueTypeName {
		case "bool":
			res[k].Value = object.Value(rk.Key).Boolean().Raw()
		case "string":
			res[k].Value = object.Value(rk.Key).String().Raw()
		case "uint":
			res[k].Value = uint(object.Value(rk.Key).Number().Raw())
		case "int":
			res[k].Value = int(object.Value(rk.Key).Number().Raw())
		case "int32":
			res[k].Value = int32(object.Value(rk.Key).Number().Raw())
		case "float64":
			res[k].Value = object.Value(rk.Key).Number().Raw()
		case "[]httptest.Responses":
			valueLen := len(res[k].Value.([]Responses))
			if rk.Length > 0 {
				valueLen = rk.Length
			}
			object.Value(rk.Key).Array().Length().IsEqual(valueLen)
			length := int(object.Value(rk.Key).Array().Length().Raw())
			if length > 0 {
				max := 1
				if rk.Length > 0 {
					max = rk.Length
				}
				if valueLen == length {
					max = length
				}
				if valueLen > 0 {
					for i := 0; i < max; i++ {
						res[k].Value.([]Responses)[i].Scan(object.Value(rk.Key).Array().Value(i).Object())
					}
				}
			}
		case "httptest.Responses":
			rk.Value.(Responses).Scan(object.Value(rk.Key).Object())
		case "[]string":
			if strings.ToLower(rk.Type) == "null" {
				res[k].Value = []string{}
			} else if strings.ToLower(rk.Type) == "notnull" {
				continue
			} else {
				length := int(object.Value(rk.Key).Array().Length().Raw())
				if length == 0 {
					continue
				}
				reskey, ok := res[k].Value.([]string)
				if ok {
					var strings []string
					for i := 0; i < length; i++ {
						strings = append(reskey, object.Value(rk.Key).Array().Value(i).String().Raw())
					}
					res[k].Value = strings
				}
			}
		default:
			continue
		}
	}
}

// Exist Check object keys if the key is in the keys array.
func Exist(object *httpexpect.Object, key string) bool {
	objectKyes := object.Keys().Raw()
	for _, objectKey := range objectKyes {
		if key == objectKey.(string) {
			return true
		}
	}
	return false
}

// GetString return string value.
func (res Responses) GetString(key ...string) string {
	if len(key) == 0 {
		return ""
	}

	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return ""
			}
			key = keys
		}
	}

	var wg sync.WaitGroup
	v := ""
	for i := 0; i < len(key); i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for m, rk := range res {
				if rk.Value == nil {
					continue
				}
				reflectTypeString := reflect.TypeOf(rk.Value).String()
				if key[i] == rk.Key {
					switch reflectTypeString {
					case "string":
						v = rk.Value.(string)
					case "httptest.Responses":
						v = res[m].Value.(Responses).GetString(key[i+1:]...)
					}
				}
			}
		}(i)
	}
	wg.Wait()
	return v
}

// GetStrArray return string array value.
func (rks Responses) GetStrArray(key string) []string {
	var wg sync.WaitGroup
	v := []string{}
	for _, rk := range rks {
		wg.Add(1)
		go func(r Response, k string) {
			defer wg.Done()
			if k == r.Key {
				if rk.Value == nil {
					return
				}
				switch reflect.TypeOf(r.Value).String() {
				case "[]string":
					v = r.Value.([]string)
				}
			}
		}(rk, key)
	}
	wg.Wait()
	return v
}

// GetResponses return Resposnes Array value
func (rks Responses) GetResponses(key string) []Responses {
	var wg sync.WaitGroup
	var v []Responses
	for _, rk := range rks {
		wg.Add(1)
		go func(r Response, k string) {
			defer wg.Done()
			if k == r.Key {
				if rk.Value == nil {
					return
				}
				switch reflect.TypeOf(r.Value).String() {
				case "[]httptest.Responses":
					v = r.Value.([]Responses)
				}
			}
		}(rk, key)
	}
	wg.Wait()
	return v
}

// GetResponsereturn Resposnes value
func (rks Responses) GetResponse(key string) Responses {
	var wg sync.WaitGroup
	var v Responses
	for _, rk := range rks {
		wg.Add(1)
		go func(r Response, k string) {
			defer wg.Done()
			if k == r.Key {
				if rk.Value == nil {
					return
				}
				switch reflect.TypeOf(r.Value).String() {
				case "httptest.Responses":
					v = r.Value.(Responses)
				}
			}
		}(rk, key)
	}
	wg.Wait()
	return v
}

// GetUint return uint value
func (rks Responses) GetUint(key ...string) uint {
	if len(key) == 0 {
		return 0
	}

	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return 0
			}
			key = keys
		}
	}

	var v uint
	var wg sync.WaitGroup
	for i := 0; i < len(key); i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for m, rk := range rks {
				if key[x] == rk.Key {
					if rk.Value == nil {
						continue
					}
					valueTypeName := reflect.TypeOf(rk.Value).String()
					switch valueTypeName {
					case "float64":
						v = uint(rk.Value.(float64))
					case "int32":
						v = uint(rk.Value.(int32))
					case "uint":
						v = rk.Value.(uint)
					case "int":
						v = uint(rk.Value.(int))
					case "httptest.Responses":
						k := key[x:]
						v = rks[m].Value.(Responses).GetUint(strings.Join(k, "."))
					}
				}
			}
		}(i)
	}
	wg.Wait()
	return v
}

// GetInt return int value
func (rks Responses) GetInt(key ...string) int {
	if len(key) == 0 {
		return 0
	}

	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return 0
			}
			key = keys
		}
	}

	var v int
	var wg sync.WaitGroup
	for i := 0; i < len(key); i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for m, rk := range rks {
				if key[x] == rk.Key {
					if rk.Value == nil {
						continue
					}
					switch reflect.TypeOf(rk.Value).String() {
					case "float64":
						v = int(rk.Value.(float64))
					case "int":
						v = rk.Value.(int)
					case "int32":
						v = int(rk.Value.(int32))
					case "uint":
						v = int(rk.Value.(uint))
					case "httptest.Responses":
						v = rks[m].Value.(Responses).GetInt(key[x+1:]...)
					}
				}
			}
		}(i)

	}
	wg.Wait()

	return v
}

// GetInt32 return int32.
func (rks Responses) GetInt32(key ...string) int32 {
	if len(key) == 0 {
		return 0
	}
	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return 0
			}
			key = keys
		}
	}
	var v int32
	var wg sync.WaitGroup
	for i := 0; i < len(key); i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for m, rk := range rks {
				if key[x] == rk.Key {
					if rk.Value == nil {
						continue
					}
					switch reflect.TypeOf(rk.Value).String() {
					case "float64":
						v = int32(rk.Value.(float64))
					case "int32":
						v = rk.Value.(int32)
					case "int":
						v = int32(rk.Value.(int))
					case "uint":
						v = int32(rk.Value.(uint))
					case "httptest.Responses":
						v = rks[m].Value.(Responses).GetInt32(key[x+1:]...)
					}
				}
			}
		}(i)
	}
	wg.Wait()

	return v
}

// GetFloat64 return float64
func (rks Responses) GetFloat64(key ...string) float64 {
	if len(key) == 0 {
		return 0
	}
	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return 0
			}
			key = keys
		}
	}
	var v float64
	var wg sync.WaitGroup
	for i := 0; i < len(key); i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for m, rk := range rks {
				if key[x] == rk.Key {
					if rk.Value == nil {
						continue
					}
					switch reflect.TypeOf(rk.Value).String() {
					case "float64":
						v = rk.Value.(float64)
					case "int":
						v = float64(rk.Value.(int))
					case "int32":
						v = float64(rk.Value.(int32))
					case "uint":
						v = float64(rk.Value.(uint))
					case "httptest.Responses":
						v = rks[m].Value.(Responses).GetFloat64(key[x+1:]...)
					}
				}
			}
		}(i)

	}
	wg.Wait()

	return v
}

// GetId return id
func (res Responses) GetId(key ...string) uint {
	if len(key) == 0 {
		key = append(key, "data", "id")
	}
	u := res.GetUint(key...)
	return u
}

var NotEmptyKey = arr.NewCheckArrayType(0)

// Schema
func Schema(str []byte) (Responses, error) {
	objs := Responses{}
	j := map[string]any{}
	if err := json.Unmarshal(str, &j); err != nil {
		return objs, fmt.Errorf("json unmarshal error %w", err)
	}
	if o, err := schema(j); err != nil {
		return objs, err
	} else {
		objs = o
	}
	return objs, nil
}

func (r Responses) Replace(key string, value any, testType ...string) {
	if len(r) == 0 {
		return
	}

	if !strings.Contains(key, ".") {
		for i1, k1 := range r {
			if k1.Key == key {
				r[i1].Value = value
				if len(testType) > 0 && testType[0] != "" {
					r[i1].Type = testType[0]
				}
			}
		}
		return
	}
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		for i1, k1 := range r {
			if k1.Key == keys[0] {
				r[i1].Value = value
				if len(testType) > 0 && testType[0] != "" {
					r[i1].Type = testType[0]
				}
			}
		}
		return
	}
	for i1, k1 := range r {
		if k1.Key != keys[0] || k1.Value == nil {
			continue
		}
		tof := reflect.TypeOf(k1.Value).String()
		if tof == "httptest.Responses" {
			r[i1].Value.(Responses).Replace(strings.Join(keys[1:], "."), value, testType...)
		} else if tof == "[]httptest.Responses" {
			if len(keys) <= 1 {
				continue
			}
			key1, _ := strconv.Atoi(keys[1])
			if r[i1].Value.([]Responses)[key1] != nil {
				r[i1].Value.([]Responses)[key1].Replace(strings.Join(keys[2:], "."), value, testType...)
			}
		}
	}
}

// schema
func schema(j map[string]any) (Responses, error) {
	objs := Responses{}
	if j == nil {
		return objs, nil
	}
	for k, v := range j {
		if k == "" {
			continue
		}
		obj := schemaResponse(k, v)
		objs = append(objs, obj)
	}
	return objs, nil
}

// schemaResponse
func schemaResponse(k string, v any) Response {
	obj := Response{}
	obj.Key = k

	if v == nil {
		return obj
	}
	typeName := reflect.TypeOf(v).String()
	switch typeName {
	case "bool":
		obj.Value = v.(bool)
	case "string":
		if obj.Key == "createdAt" || obj.Key == "updatedAt" || obj.Key == "deletedAt" {
			obj.Type = "notempty"
		} else if NotEmptyKey.Len() > 0 && NotEmptyKey.Check(obj.Key) {
			obj.Type = "notempty"
		} else {
			obj.Value = v.(string)
		}
	case "uint":
		obj.Value = v.(uint)
	case "int":
		obj.Value = v.(int)
	case "int32":
		obj.Value = v.(int32)
	case "float64":
		obj.Value = v.(float64)
	case "[]string":
		obj.Value = v.([]string)
	case "map[string]interface {}":
		if value, _ := schema(v.(map[string]any)); value != nil {
			obj.Value = value
		}
	case "[]interface {}":
		list := []Responses{}
		for _, v1 := range v.([]any) {
			listObj := Responses{}
			if v3, ok := v1.(map[string]any); ok {
				for k2, v2 := range v3 {
					listObj = append(listObj, schemaResponse(k2, v2))
				}
				list = append(list, listObj)
				obj.Value = list
			} else if _, ok := v1.(string); ok {
				obj.Value = v
			}
		}

	default:
		fmt.Printf("schemaResponse key:%s valueTypeName:%s\n", k, typeName)
	}
	return obj
}
