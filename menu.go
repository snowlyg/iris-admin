package admin

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Path       string `json:"path"`
	Component  string `json:"component"`
	Redirect   string `json:"redirect"`
	Hidden     bool   `json:"hidden"`
	AlwaysShow bool   `json:"alwaysShow"`
	Meta
	Children []*Menu `json:"children" gorm:"-"`
}

func (m *Menu) TableName() string {
	return "menus"
}

type Meta struct {
	Roles   []string `json:"roles" gorm:"-"`
	Title   string   `json:"title"`
	Icon    string   `json:"icon"`
	NoCache bool     `json:"noCache"`
}
