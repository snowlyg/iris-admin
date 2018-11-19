package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PlanDeparts struct {
	Id        int       `orm:"column(id);auto"`
	Name      string    `orm:"column(name);size(191)" description:"库名部门名称"`
	Icon      string    `orm:"column(icon);size(191);null" description:"库名部门图标"`
	PlanId    *Plans    `orm:"column(plan_id);rel(fk)"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt time.Time `orm:"column(updated_at);type(timestamp);null"`
	DeletedAt time.Time `orm:"column(deleted_at);type(timestamp);null"`
}

func (t *PlanDeparts) TableName() string {
	return "plan_departs"
}

func init() {
	orm.RegisterModel(new(PlanDeparts))
}

// AddPlanDeparts insert a new PlanDeparts into database and returns
// last inserted Id on success.
func AddPlanDeparts(m *PlanDeparts) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPlanDepartsById retrieves PlanDeparts by Id. Returns error if
// Id doesn't exist
func GetPlanDepartsById(id int) (v *PlanDeparts, err error) {
	o := orm.NewOrm()
	v = &PlanDeparts{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPlanDeparts retrieves all PlanDeparts matches certain condition. Returns empty list if
// no records exist
func GetAllPlanDeparts(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PlanDeparts))
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

	var l []PlanDeparts
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

// UpdatePlanDeparts updates PlanDeparts by Id and returns error if
// the record to be updated doesn't exist
func UpdatePlanDepartsById(m *PlanDeparts) (err error) {
	o := orm.NewOrm()
	v := PlanDeparts{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePlanDeparts deletes PlanDeparts by Id and returns error if
// the record to be deleted doesn't exist
func DeletePlanDeparts(id int) (err error) {
	o := orm.NewOrm()
	v := PlanDeparts{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PlanDeparts{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
