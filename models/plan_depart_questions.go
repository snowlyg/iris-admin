package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PlanDepartQuestions struct {
	Id                 int          `orm:"column(id);auto"`
	Question           string       `orm:"column(question);size(191)" description:"问题"`
	Answers            string       `orm:"column(answers);size(191)" description:"答案"`
	Status             string       `orm:"column(status);size(191)" description:"问题状态 ‘未开始’，‘调查中’，‘审核中’，‘审核成功’，‘已驳回’"`
	ClientAnswer       string       `orm:"column(client_answer);size(191);null" description:"客户回答"`
	ClientAnswerEditer string       `orm:"column(client_answer_editer);size(191);null" description:"调查过程执行人"`
	Auditer            string       `orm:"column(auditer);size(191);null" description:"问题审核人"`
	AuditText          string       `orm:"column(audit_text);size(191);null" description:"审核备注"`
	AuditAt            string       `orm:"column(audit_at);size(191);null" description:"审核时间"`
	MoreFiles          string       `orm:"column(more_files);null" description:"补充材料 图片url 数组"`
	ConfirmText        string       `orm:"column(confirm_text);null" description:"现场确认内容 需要富文本"`
	ConfirmEditer      string       `orm:"column(confirm_editer);size(191);null" description:"现场确认执行人"`
	ConfirmAt          time.Time    `orm:"column(confirm_at);type(datetime);null" description:"现场确认内容编辑时间"`
	ConclusionStatus   string       `orm:"column(conclusion_status);size(191);null" description:"最终结论严重性"`
	ConclusionAt       time.Time    `orm:"column(conclusion_at);type(datetime);null" description:"最终结论 编辑时间"`
	Conclusion         string       `orm:"column(conclusion);null" description:"最终结论 需要富文本"`
	ConclusionEditer   string       `orm:"column(conclusion_editer);size(191);null" description:"最终结论执行人"`
	PlanDepartId       *PlanDeparts `orm:"column(plan_depart_id);rel(fk)"`
	CreatedAt          time.Time    `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt          time.Time    `orm:"column(updated_at);type(timestamp);null"`
	DeletedAt          time.Time    `orm:"column(deleted_at);type(timestamp);null"`
}

func (t *PlanDepartQuestions) TableName() string {
	return "plan_depart_questions"
}

func init() {
	orm.RegisterModel(new(PlanDepartQuestions))
}

// AddPlanDepartQuestions insert a new PlanDepartQuestions into database and returns
// last inserted Id on success.
func AddPlanDepartQuestions(m *PlanDepartQuestions) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPlanDepartQuestionsById retrieves PlanDepartQuestions by Id. Returns error if
// Id doesn't exist
func GetPlanDepartQuestionsById(id int) (v *PlanDepartQuestions, err error) {
	o := orm.NewOrm()
	v = &PlanDepartQuestions{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPlanDepartQuestions retrieves all PlanDepartQuestions matches certain condition. Returns empty list if
// no records exist
func GetAllPlanDepartQuestions(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PlanDepartQuestions))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []PlanDepartQuestions
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdatePlanDepartQuestions updates PlanDepartQuestions by Id and returns error if
// the record to be updated doesn't exist
func UpdatePlanDepartQuestionsById(m *PlanDepartQuestions) (err error) {
	o := orm.NewOrm()
	v := PlanDepartQuestions{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePlanDepartQuestions deletes PlanDepartQuestions by Id and returns error if
// the record to be deleted doesn't exist
func DeletePlanDepartQuestions(id int) (err error) {
	o := orm.NewOrm()
	v := PlanDepartQuestions{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PlanDepartQuestions{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
