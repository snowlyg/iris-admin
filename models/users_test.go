package models

import (
	"reflect"
	"testing"

	"IrisAdminApi/transformer"
	"IrisAdminApi/validates"
	"github.com/jinzhu/gorm"
)

func TestGetAllUsers(t *testing.T) {
	type args struct {
		name    string
		orderBy string
		offset  int
		limit   int
	}
	tests := []struct {
		name string
		args args
		want []*User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllUsers(tt.args.name, tt.args.orderBy, tt.args.offset, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	type args struct {
		pwd string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashPassword(tt.args.pwd); got != tt.want {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	type args struct {
		id       uint
		username string
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUser(tt.args.id, tt.args.username); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAdminLogout(t *testing.T) {
	type args struct {
		userId uint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserAdminLogout(tt.args.userId); got != tt.want {
				t.Errorf("UserAdminLogout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_CheckLogin(t *testing.T) {
	type fields struct {
		Model    gorm.Model
		Name     string
		Username string
		Password string
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Token
		want1  bool
		want2  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Model:    tt.fields.Model,
				Name:     tt.fields.Name,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
			got, got1, got2 := u.CheckLogin(tt.args.password)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckLogin() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckLogin() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("CheckLogin() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestUser_CreateSystemAdmin(t *testing.T) {
	type fields struct {
		Model    gorm.Model
		Name     string
		Username string
		Password string
	}
	type args struct {
		roleId uint
		rc     *transformer.Conf
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Model:    tt.fields.Model,
				Name:     tt.fields.Name,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
		})
	}
}

func TestUser_CreateUser(t *testing.T) {
	type fields struct {
		Model    gorm.Model
		Name     string
		Username string
		Password string
	}
	type args struct {
		aul *validates.CreateUpdateUserRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Model:    tt.fields.Model,
				Name:     tt.fields.Name,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
		})
	}
}

func TestUser_DeleteUser(t *testing.T) {
	type fields struct {
		Model    gorm.Model
		Name     string
		Username string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Model:    tt.fields.Model,
				Name:     tt.fields.Name,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
		})
	}
}

func TestUser_GetUserById(t *testing.T) {
	type fields struct {
		Model    gorm.Model
		Name     string
		Username string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Model:    tt.fields.Model,
				Name:     tt.fields.Name,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
		})
	}
}

func TestUser_GetUserByUsername(t *testing.T) {
	type fields struct {
		Model    gorm.Model
		Name     string
		Username string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Model:    tt.fields.Model,
				Name:     tt.fields.Name,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
		})
	}
}

func TestUser_UpdateUser(t *testing.T) {
	type fields struct {
		Model    gorm.Model
		Name     string
		Username string
		Password string
	}
	type args struct {
		uj *validates.CreateUpdateUserRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Model:    tt.fields.Model,
				Name:     tt.fields.Name,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
			}
		})
	}
}

func Test_addRoles(t *testing.T) {
	type args struct {
		uj   *validates.CreateUpdateUserRequest
		user *User
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
