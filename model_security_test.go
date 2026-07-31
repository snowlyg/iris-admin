package admin

import "testing"

func TestGetOrderBySanitizesInput(t *testing.T) {
	tests := []struct {
		name    string
		sort    string
		orderBy string
		want    string
	}{
		{name: "defaults", want: "created_at desc"},
		{name: "valid", sort: "ASC", orderBy: "users.created_at", want: "users.created_at asc"},
		{name: "invalid direction", sort: "desc nulls last", orderBy: "id", want: "id desc"},
		{name: "invalid expression", sort: "asc", orderBy: "created_at desc; DROP TABLE users", want: "created_at asc"},
		{name: "invalid function", sort: "asc", orderBy: "RAND()", want: "created_at asc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOrderBy(tt.sort, tt.orderBy); got != tt.want {
				t.Fatalf("getOrderBy(%q, %q) = %q, want %q", tt.sort, tt.orderBy, got, tt.want)
			}
		})
	}
}
