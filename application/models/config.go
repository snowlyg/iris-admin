package models

import (
	"gorm.io/gorm"
)

type Config struct {
	gorm.Model
	Name  string `gorm:"unique;not null;type:varchar(256)" json:"name" validate:"required" comment:"name"`
	Value string `gorm:"not null;type:varchar(1024)" json:"value" validate:"required"  comment:"value"`
}
