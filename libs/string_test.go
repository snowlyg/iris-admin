// +build test

package libs

import (
	"testing"
)

func TestBase64Decode(t *testing.T) {
	s := "MTIzNDU2"
	want := "123456"
	t.Run("TestBase64Decode", func(t *testing.T) {
		if got := Base64Decode(s); got != want {
			t.Errorf("Base64Decode() = %v, want %v", got, want)
		}
	})
}

func TestHashPassword(t *testing.T) {
	notWant := ""
	t.Run("TestHashPassword", func(t *testing.T) {
		if got := HashPassword("TestHashPassword"); got == notWant {
			t.Errorf("HashPassword() = %v, not_want %v", got, notWant)
		}
	})

}

func TestParseFlostToString(t *testing.T) {
	args := 3.140
	want := "3.14000"
	t.Run("TestParseFlostToString", func(t *testing.T) {
		if got := ParseFlostToString(args); got != want {
			t.Errorf("ParseFlostToString() = %v, want %v", got, want)
		}
	})
}

func TestParseInt(t *testing.T) {
	type args struct {
		b      string
		defInt int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "string转换int",
			args: args{
				b:      "314",
				defInt: 315,
			},
			want: 314,
		},
		{
			name: "string转换int 得到默认值",
			args: args{
				b:      "",
				defInt: 315,
			},
			want: 315,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseInt(tt.args.b, tt.args.defInt); got != tt.want {
				t.Errorf("ParseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	type args struct {
		b int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "int转换string",
			args: args{314},
			want: "314",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseString(tt.args.b); got != tt.want {
				t.Errorf("ParseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubString(t *testing.T) {
	type args struct {
		str    string
		start  int
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "字符串截取全部",
			args: args{
				str:    "TestSubString",
				start:  0,
				length: 0,
			},
			want: "",
		},
		{
			name: "字符串截取 start 小于 0",
			args: args{
				str:    "TestSubString",
				start:  -1,
				length: 1,
			},
			want: "g",
		},
		{
			name: "字符串截取 start 小于 0",
			args: args{
				str:    "TestSubString",
				start:  len("TestSubString") + 1,
				length: 1,
			},
			want: "",
		},
		{
			name: "字符串截取 start 大于 字符长度",
			args: args{
				str:    "TestSubString",
				start:  len("TestSubString") + 1,
				length: 1,
			},
			want: "",
		},
		{
			name: "字符串截取 length 小于 0",
			args: args{
				str:    "TestSubString",
				start:  len("TestSubString") + 1,
				length: -1,
			},
			want: "g",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubString(tt.args.str, tt.args.start, tt.args.length); got != tt.want {
				t.Errorf("SubString() = %v, want %v", got, tt.want)
			}
		})
	}
}
