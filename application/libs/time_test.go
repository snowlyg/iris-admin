// +build test

package libs

import (
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	location := &time.Location{}
	tt := time.Date(2020, 01, 02, 15, 04, 05, 0, location)
	want := "2020-01-02 15:04:05"
	t.Run("TestTimeFormat", func(t *testing.T) {
		if got := TimeFormat(&tt, ""); got != want {
			t.Errorf("TimeFormat() = %v, want %v", got, want)
		}
	})
}
