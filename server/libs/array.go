package libs

import (
	"strconv"
)

//InArrayS 如果 s 在 items 中,返回 true；否则，返回 false。
func InArrayS(items []string, s string) bool {
	for _, item := range items {
		if item == s {
			return true
		}
	}
	return false
}

// 连接 unit slice 为字符串
func UnitJoin(ss []uint, sep string) string {
	var rs string
	for index, item := range ss {
		itemS := strconv.FormatUint(uint64(item), 10)
		if index < len(ss)-1 {
			rs += itemS + sep
		} else {
			rs += itemS
		}
	}
	return rs
}
