// +build test public api tag access perm role user type doc chapter article expire config article

package mock

import (
	"github.com/bxcodec/faker/v3"
	"reflect"
	"time"
)

// CustomGenerator ...
func CustomGenerator() {
	_ = faker.AddProvider("ctime", func(v reflect.Value) (interface{}, error) {
		return time.Now().Format("2006-01-02T15:04:05Z"), nil
	})

}

type Article struct {
	Title        string `faker:"len=200,lang=chi" json:"title"`
	ContentShort string `faker:"len=200,lang=chi" json:"content_short"`
	Author       string `faker:"len=25,lang=chi" json:"author"`
	ImageUri     string `faker:"url" json:"image_uri"`
	SourceUri    string `faker:"url" json:"source_uri"`
	IsOriginal   bool
	Content      string `faker:"len=1000,lang=chi" json:"content"`
	Status       string `faker:"len=10,lang=chi" json:"status"`
	DisplayTime  string `faker:"ctime" json:"display_time"`
	Like         int64  `json:"like"`
	Read         int64  `json:"read"`
	Ips          string `faker:"paragraph,lang=chi" json:"ips"`
}

type Chapter struct {
	Title        string `faker:"len=200,lang=chi" json:"title"`
	ContentShort string `faker:"len=200,lang=chi" json:"content_short"`
	Author       string `faker:"len=20,lang=chi" json:"author"`
	ImageUri     string `faker:"url" json:"image_uri"`
	SourceUri    string `faker:"url" json:"source_uri"`
	IsOriginal   bool   `json:"is_original"`
	Content      string `faker:"len=2000,lang=chi" json:"content"`
	Status       string `faker:"len=10,lang=chi" json:"status"`
	DisplayTime  string `faker:"ctime" json:"display_time"`
	Like         int64  `json:"like"`
	Read         int64  `json:"read"`
	Ips          string `faker:"len=200,lang=chi" json:"ips"`
	Sort         int64  `faker:"boundary_start=31, boundary_end=88" json:"sort"`
}

type Config struct {
	Name  string `faker:"len=25,unique" json:"name"`
	Value string `faker:"len=25,lang=chi" json:"value"`
}

type Doc struct {
	Name string `faker:"len=25,lang=chi" json:"name"`
}

type Permission struct {
	Name        string `faker:"len=25,unique" json:"name"`
	DisplayName string `faker:"len=25,lang=chi" json:"display_name"`
	Description string `faker:"len=25,lang=chi" json:"description"`
	Act         string `faker:"len=25" json:"act"`
}

type Role struct {
	Name        string `faker:"len=25,unique" json:"name"`
	DisplayName string `faker:"len=25,lang=chi" json:"display_name"`
	Description string `faker:"len=25,lang=chi" json:"description"`
}

type Tag struct {
	Name string `faker:"len=25,lang=chi" json:"name"`
}

type Type struct {
	Name string `faker:"len=25,lang=chi" json:"name"`
}

type User struct {
	Name     string `faker:"len=25,lang=chi" json:"name"`
	Username string `faker:"len=25,unique" json:"username"`
	Password string `faker:"password" json:"password"`
	Intro    string `faker:"len=25,lang=chi" json:"intro"`
	Avatar   string `faker:"url" json:"avatar"`
	RoleIds  []uint `json:"role_ids"`
}
