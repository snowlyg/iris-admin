package libs

import "testing"

func TestInArrayS(t *testing.T) {
	type args struct {
		items []string
		s     string
	}
	items := []string{"1", "error", ";"}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: struct {
				items []string
				s     string
			}{items: items, s: "1"},
			want: true,
		},
		{
			name: "success 2",
			args: struct {
				items []string
				s     string
			}{items: items, s: "error"},
			want: true,
		},
		{
			name: "success 3",
			args: struct {
				items []string
				s     string
			}{items: items, s: ";"},
			want: true,
		},
		{
			name: "error",
			args: struct {
				items []string
				s     string
			}{items: items, s: "2"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InArrayS(tt.args.items, tt.args.s); got != tt.want {
				t.Errorf("InArrayS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnitJoin(t *testing.T) {
	type args struct {
		ss  []uint
		sep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: struct {
				ss  []uint
				sep string
			}{ss: []uint{1, 2, 3, 4}, sep: ","},
			want: "1,2,3,4",
		},
		{
			name: "success",
			args: struct {
				ss  []uint
				sep string
			}{ss: []uint{1, 2, 3, 4}, sep: "||"},
			want: "1||2||3||4",
		},
		{
			name: "success",
			args: struct {
				ss  []uint
				sep string
			}{ss: []uint{1, 2, 3, 4}, sep: "-"},
			want: "1-2-3-4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnitJoin(tt.args.ss, tt.args.sep); got != tt.want {
				t.Errorf("UnitJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}
