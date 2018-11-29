package main

import (
	"github.com/jinzhu/gorm"
)

type QuestionDetails struct {
	gorm.Model
	Question             string `gorm:"not null comment('追加问题') VARCHAR(191)"`
	Answer               string `gorm:"not null comment('追加问题回答') LONGTEXT"`
	PlanDepartQuestionId int    `gorm:"not null index INT(10)"`
}
