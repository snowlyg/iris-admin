// +build test

package libs

import "testing"

func TestApiResource(t *testing.T) {
	type args struct {
		code    int64
		objects interface{}
		msg     string
	}

	error403 := args{code: 403, objects: "1", msg: "403 error"}
	error401 := args{code: 401, objects: "1", msg: "401 error"}
	error400 := args{code: 400, objects: "1", msg: "400 error"}
	success200 := args{code: 200, objects: "1", msg: "200 error"}

	tests := []struct {
		name string
		args args
		want args
	}{
		{
			name: "400 error",
			args: error403,
			want: error403,
		},
		{
			name: "403 error",
			args: error401,
			want: error401,
		},
		{
			name: "401 error",
			args: error400,
			want: error400,
		},
		{
			name: "200 success",
			args: success200,
			want: success200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApiResource(tt.args.code, tt.args.objects, tt.args.msg)
			if got.Code != tt.want.code || got.Msg != tt.want.msg || got.Data != tt.want.objects {
				t.Errorf("TestApiResource() = %+v\n, want %+v\n", got, tt.want)
			}
		})
	}
}
