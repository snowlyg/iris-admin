package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ReportDepartQuestions struct {
	gorm.Model
	Question           string    `gorm:"not null comment('问题') VARCHAR(191)"`
	Answers            string    `gorm:"not null comment('答案') VARCHAR(191)"`
	Status             string    `gorm:"not null default '未开始' comment('问题状态 ‘未开始’，‘调查中’，‘审核中’，‘审核成功’，‘已驳回’') VARCHAR(191)"`
	ClientAnswer       string    `gorm:"comment('客户回答') VARCHAR(191)"`
	ClientAnswerEditer string    `gorm:"comment('调查过程执行人') VARCHAR(191)"`
	Auditer            string    `gorm:"comment('问题审核人') VARCHAR(191)"`
	AuditText          string    `gorm:"comment('审核备注') VARCHAR(191)"`
	AuditAt            string    `gorm:"comment('审核时间') VARCHAR(191)"`
	MoreFiles          string    `gorm:"comment('补充材料 图片url 数组') LONGTEXT"`
	ConfirmText        string    `gorm:"comment('现场确认内容 需要富文本') LONGTEXT"`
	ConfirmEditer      string    `gorm:"comment('现场确认执行人') VARCHAR(191)"`
	ConfirmAt          time.Time `gorm:"comment('现场确认内容编辑时间') DATETIME"`
	ConclusionStatus   string    `gorm:"comment('最终结论严重性') VARCHAR(191)"`
	ConclusionAt       time.Time `gorm:"comment('最终结论 编辑时间') DATETIME"`
	Conclusion         string    `gorm:"comment('最终结论 需要富文本') LONGTEXT"`
	ConclusionEditer   string    `gorm:"comment('最终结论执行人') VARCHAR(191)"`
	QuestionDetails    string    `gorm:"comment('进一步提问 array') LONGTEXT"`
	ReportDepartId     int       `gorm:"not null index INT(10)"`
}
