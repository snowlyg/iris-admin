package models

import (
	"time"
)

type QuestionDetails struct {
	Id                   int       `xorm:"not null pk autoincr INT(10)"`
	Question             string    `xorm:"not null comment('追加问题') VARCHAR(191)"`
	Answer               string    `xorm:"not null comment('追加问题回答') LONGTEXT"`
	PlanDepartQuestionId int       `xorm:"not null index INT(10)"`
	CreatedAt            time.Time `xorm:"TIMESTAMP"`
	UpdatedAt            time.Time `xorm:"TIMESTAMP"`
	DeletedAt            time.Time `xorm:"TIMESTAMP"`
}
