package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Revisions struct {
	Id               int       `orm:"column(id);auto"`
	RevisionableType string    `orm:"column(revisionable_type);size(191)"`
	RevisionableId   int       `orm:"column(revisionable_id)"`
	UserId           int       `orm:"column(user_id);null"`
	Key              string    `orm:"column(key);size(191)"`
	OldValue         string    `orm:"column(old_value);null"`
	NewValue         string    `orm:"column(new_value);null"`
	Device           string    `orm:"column(device);size(191);null"`
	Ip               string    `orm:"column(ip);size(191);null"`
	DeviceType       string    `orm:"column(device_type);size(191);null"`
	Address          string    `orm:"column(address);size(191);null"`
	Browser          string    `orm:"column(browser);size(191);null"`
	Platform         string    `orm:"column(platform);size(191);null"`
	Language         string    `orm:"column(language);size(191);null"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);null"`
}

func (t *Revisions) TableName() string {
	return "revisions"
}

func init() {
	orm.RegisterModel(new(Revisions))
}

// AddRevisions insert a new Revisions into database and returns
// last inserted Id on success.
func AddRevisions(m *Revisions) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRevisionsById retrieves Revisions by Id. Returns error if
// Id doesn't exist
func GetRevisionsById(id int) (v *Revisions, err error) {
	o := orm.NewOrm()
	v = &Revisions{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRevisions retrieves all Revisions matches certain condition. Returns empty list if
// no records exist
func GetAllRevisions(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Revisions))
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

	var l []Revisions
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

// UpdateRevisions updates Revisions by Id and returns error if
// the record to be updated doesn't exist
func UpdateRevisionsById(m *Revisions) (err error) {
	o := orm.NewOrm()
	v := Revisions{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRevisions deletes Revisions by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRevisions(id int) (err error) {
	o := orm.NewOrm()
	v := Revisions{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Revisions{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
