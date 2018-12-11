package tools

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Tools struct {
}

var (
	Tool = New()
	once sync.Once
)

/**
 * 返回单例实例
 * @method New
 */
func New() (t *Tools) {
	once.Do(func() { //只执行一次
		t = &Tools{}
	})
	return t
}

/**
 * string转换int
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) ParseInt(b string, defInt int) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		return defInt
	} else {
		return id
	}
}

/**
 * int转换string
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) ParseString(b int) string {
	id := strconv.Itoa(b)
	return id
}

/**
 * 转换浮点数为string
 * @method func
 * @param  {[type]} t *             Tools [description]
 * @return {[type]}   [description]
 */
func (t *Tools) ParseFlostToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 5, 64)
}

/**
 * 结构体转成json 字符串
 * @method StruckToString
 * @param  {[type]}       data interface{} [description]
 */
func (t *Tools) StructToString(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

/**
 * 结构体转换成map对象
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) StructToMap(obj interface{}) map[string]interface{} {
	k := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < k.NumField(); i++ {
		data[strings.ToLower(k.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

//生成随机字符串
func (t *Tools) GetRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()+[]{}/<>;:=.,?"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

/**
 * 字符串截取
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) SubString(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	rune_str := []rune(str)
	len_str := len(rune_str)

	if start < 0 {
		start = len_str + start
	}
	if start > len_str {
		start = len_str
	}
	end := start + length
	if end > len_str {
		end = len_str
	}
	if length < 0 {
		end = len_str + length
	}
	if start > end {
		start, end = end, start
	}
	return string(rune_str[start:end])
}

/**
 * base64 解码
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) Base64Decode(str string) string {
	s, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(s)
}

/**
 * base64 解码
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) TimeFormat(time *time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
