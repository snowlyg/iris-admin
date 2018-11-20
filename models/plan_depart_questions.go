package models

import (
	"time"
)

type PlanDepartQuestions struct {
	Id                 int       `xorm:"not null pk autoincr INT(10)"`
	Question           string    `xorm:"not null comment('问题') VARCHAR(191)"`
	Answers            string    `xorm:"not null comment('答案') VARCHAR(191)"`
	Status             string    `xorm:"not null default '未开始' comment('问题状态 ‘未开始’，‘调查中’，‘审核中’，‘审核成功’，‘已驳回’') VARCHAR(191)"`
	ClientAnswer       string    `xorm:"comment('客户回答') VARCHAR(191)"`
	ClientAnswerEditer string    `xorm:"comment('调查过程执行人') VARCHAR(191)"`
	Auditer            string    `xorm:"comment('问题审核人') VARCHAR(191)"`
	AuditText          string    `xorm:"comment('审核备注') VARCHAR(191)"`
	AuditAt            string    `xorm:"comment('审核时间') VARCHAR(191)"`
	MoreFiles          string    `xorm:"comment('补充材料 图片url 数组') LONGTEXT"`
	ConfirmText        string    `xorm:"comment('现场确认内容 需要富文本') LONGTEXT"`
	ConfirmEditer      string    `xorm:"comment('现场确认执行人') VARCHAR(191)"`
	ConfirmAt          time.Time `xorm:"comment('现场确认内容编辑时间') DATETIME"`
	ConclusionStatus   string    `xorm:"comment('最终结论严重性') VARCHAR(191)"`
	ConclusionAt       time.Time `xorm:"comment('最终结论 编辑时间') DATETIME"`
	Conclusion         string    `xorm:"comment('最终结论 需要富文本') LONGTEXT"`
	ConclusionEditer   string    `xorm:"comment('最终结论执行人') VARCHAR(191)"`
	PlanDepartId       int       `xorm:"not null index INT(10)"`
	CreatedAt          time.Time `xorm:"TIMESTAMP"`
	UpdatedAt          time.Time `xorm:"TIMESTAMP"`
	DeletedAt          time.Time `xorm:"TIMESTAMP"`
}
