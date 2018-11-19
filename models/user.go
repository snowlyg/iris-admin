package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"strings"
	"time"
)

type Users struct {
	Id               int       `orm:"column(id);auto"`
	Name             string    `orm:"column(name);size(191)"`
	Username         string    `orm:"column(username);size(191);null"`
	Password         string    `orm:"column(password);size(191)"`
	Confirmed        int8      `orm:"column(confirmed)"`
	IsClient         int8      `orm:"column(is_client)"`
	IsFrozen         int8      `orm:"column(is_frozen)"`
	IsAudit          int8      `orm:"column(is_audit)"`
	IsClientAdmin    int8      `orm:"column(is_client_admin)"`
	WechatName       string    `orm:"column(wechat_name);size(191);null"`
	WechatAvatar     string    `orm:"column(wechat_avatar);size(191);null"`
	Email            string    `orm:"column(email);size(191);null"`
	OpenId           string    `orm:"column(open_id);size(191);null"`
	WechatVerfiyTime time.Time `orm:"column(wechat_verfiy_time);type(datetime);null"`
	IsWechatVerfiy   int8      `orm:"column(is_wechat_verfiy)"`
	Phone            string    `orm:"column(phone);size(191);null"`
	RememberToken    string    `orm:"column(remember_token);size(100);null"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);null"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);null"`
	DeletedAt        time.Time `orm:"column(deleted_at);type(timestamp);null"`
}

func (t *Users) TableName() string {
	return "users"
}

func init() {

	orm.RegisterModel(new(Users))
}

// AddUsers insert a new Users into database and returns
// last inserted Id on success.
func AddUsers(m *Users) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUsersById retrieves Users by Id. Returns error if
// Id doesn't exist
func GetUsersById(id int) (v *Users, err error) {
	o := orm.NewOrm()
	o.Using("default")

	v = &Users{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUsers retrieves all Users matches certain condition. Returns empty list if
// no records exist
func GetAllUsers(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Users))
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

	var l []Users
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

// UpdateUsers updates Users by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsersById(m *Users) (err error) {
	o := orm.NewOrm()
	v := Users{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUsers deletes Users by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUsers(id int) (err error) {
	o := orm.NewOrm()
	v := Users{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Users{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// user's input.
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword will check if passwords are matched.
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}

//验证帐号密码
func UserAdminCheckLogin(userlogin *UserLogin) (user *Users, err error) {

	user = new(Users)

	if userlogin.Username == "" || userlogin.Password == "" {
		return user, errors.New("please fill both username and password fields")
	}

	//fmt.Println("userLogin:", userlogin)

	o := orm.NewOrm()
	user.Username = userlogin.Username

	err = o.Read(user,"username")


	if err != nil {
		return user, err
	}

	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
	//fmt.Println(string(hashedPassword))

	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(userlogin.Password))

	if err != nil {
		return user, err
	}

	return user, nil

}
