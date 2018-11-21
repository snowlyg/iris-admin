package models

import "time"

type Roles struct {
	Id          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Name        string     `gorm:"not null VARCHAR(191)"`
	GuardName   string     `gorm:"not null VARCHAR(191)"`
	DisplayName string     `gorm:"VARCHAR(191)"`
	Description string     `gorm:"VARCHAR(191)"`
	Level       int        `gorm:"not null default 0 INT(10)"`
}
