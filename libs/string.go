package libs

import (
	"encoding/base64"
	"math/rand"
	"strconv"

	"github.com/jameskeane/bcrypt"
)

/**
 * string转换int
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func ParseInt(b string, defInt int) int {
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
func ParseString(b int) string {
	id := strconv.Itoa(b)
	return id
}

/**
 * 转换浮点数为string
 * @method func
 * @param  {[type]} t *             Tools [description]
 * @return {[type]}   [description]
 */
func ParseFlostToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 5, 64)
}

//生成随机字符串
func GetRandomString(n int) string {
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
func SubString(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	runeStr := []rune(str)
	lenStr := len(runeStr)

	if start < 0 {
		start = lenStr + start
	}
	if start > lenStr {
		start = lenStr
	}
	end := start + length
	if end > lenStr {
		end = lenStr
	}
	if length < 0 {
		end = lenStr + length
	}
	if start > end {
		start, end = end, start
	}
	return string(runeStr[start:end])
}

/**
 * base64 解码
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func Base64Decode(str string) string {
	s, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(s)
}

func HashPassword(pwd string) string {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(pwd, salt)
	return hash
}

