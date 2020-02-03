package config

import (
	"testing"
)

func TestGetAppCreateSysData(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "config",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := GetAppCreateSysData(); got != tt.want {
				t.Errorf("GetAppCreateSysData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppDirverType(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "Sqlite",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppDirverType(); got != tt.want {
				t.Errorf("GetAppDirverType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppLoggerLevel(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "debug",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppLoggerLevel(); got != tt.want {
				t.Errorf("GetAppLoggerLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "IrisAdminApi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppName(); got != tt.want {
				t.Errorf("GetAppName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppURl(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "irisadminapi.com:80",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppUrl(); got != tt.want {
				t.Errorf("GetAppURl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMongodbConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "mongodb://root:123456@127.0.0.1:27017/admin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMongodbConnect(); got != tt.want {
				t.Errorf("GetMongodbConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMysqlConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "root:passwrod@(127.0.0.1:3306)/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMysqlConnect(); got != tt.want {
				t.Errorf("GetMysqlConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMysqlName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "iris",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMysqlName(); got != tt.want {
				t.Errorf("GetMysqlName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMysqlTName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "tiris",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMysqlTName(); got != tt.want {
				t.Errorf("GetMysqlTName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSqliteConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "./tmp/gorm.db",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSqliteConnect(); got != tt.want {
				t.Errorf("GetSqliteConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSqliteTConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "./tmp/tgorm.db",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSqliteTConnect(); got != tt.want {
				t.Errorf("GetSqliteTConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTestDataName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "超级管理员",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTestDataName(); got != tt.want {
				t.Errorf("GetTestDataName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTestDataPwd(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTestDataPwd(); got != tt.want {
				t.Errorf("GetTestDataPwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTestDataUserName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "username",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTestDataUserName(); got != tt.want {
				t.Errorf("GetTestDataUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}
