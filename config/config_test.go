package config

import (
	"testing"
)

func TestSetAppCreateSysData(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "config",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAppCreateSysData(false)
			if got := GetAppCreateSysData(); got != tt.want {
				t.Errorf("SetAppCreateSysData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAppDriverType(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "Mysql",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = SetAppDriverType("Mysql")
			if got := GetAppDriverType(); got != tt.want {
				t.Errorf("SetAppDriverType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAppLoggerLevel(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "info",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAppLoggerLevel("info")
			if got := GetAppLoggerLevel(); got != tt.want {
				t.Errorf("SetAppLoggerLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAppName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "SetAppName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAppName("SetAppName")
			if got := GetAppName(); got != tt.want {
				t.Errorf("SetAppName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAppURl(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAppUrl("localhost:8080")
			if got := GetAppUrl(); got != tt.want {
				t.Errorf("SetAppURl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetMongodbConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "mongodb://root:123456@127.0.0.1:27017/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetMongodbConnect("mongodb://root:123456@127.0.0.1:27017/test")
			if got := GetMongodbConnect(); got != tt.want {
				t.Errorf("SetMongodbConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetMysqlConnect(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "root:([a-z]+)@(127.0.0.1:3306)/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetMysqlConnect("root:([a-z]+)@(127.0.0.1:3306)/")
			//matchString, err := regexp.MatchString(tt.want, got)
			if got := GetMysqlConnect(); got != tt.want {
				t.Errorf("SetMysqlConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetMysqlName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "test_iris",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetMysqlName("test_iris")
			if got := GetMysqlName(); got != tt.want {
				t.Errorf("SetMysqlName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetMysqlTName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "ttiris",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetMysqlTName("ttiris")
			if got := GetMysqlTName(); got != tt.want {
				t.Errorf("SetMysqlTName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetTestDataName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "管理员",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTestDataName("管理员")
			if got := GetTestDataName(); got != tt.want {
				t.Errorf("SetTestDataName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetTestDataPwd(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "pwd123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTestDataPwd("pwd123456")
			if got := GetTestDataPwd(); got != tt.want {
				t.Errorf("GetTestDataPwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetTestDataUserName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "config",
			want: "username123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTestDataUserName("username123456")
			if got := GetTestDataUserName(); got != tt.want {
				t.Errorf("GetTestDataUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}
