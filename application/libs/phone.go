package libs

import (
	"fmt"
	"math/rand"
	"time"
)

var prelist = []string{"130", "131", "132", "133", "134", "135", "136", "137", "138", "139", "147", "150", "151", "152", "153", "155", "156", "157", "158", "159", "181", "182", "183", "184", "185", "186", "187", "188", "189", "191"}

type GeneratePhoneNumber struct {
	CacheData []string
}

// 生成随机手机号码
func (*GeneratePhoneNumber) CreatePhoneNumber() string {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return prelist[rd.Int63n(28)] + fmt.Sprintf("%08v", rd.Int63n(100000000))
}

// 生成唯一随机手机号码
func (g *GeneratePhoneNumber) CreateUniquePhoneNumber() string {
	var pn string
	for true {
		pn = g.CreatePhoneNumber()
		if !InArrayS(g.CacheData, pn) {
			g.CacheData = append(g.CacheData, pn)
			return pn
		}
	}
	return pn
}
