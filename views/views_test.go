package views

import (
	"testing"
	"time"
)

func TestNiceDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17-Dec-20 04:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17-Dec-20 03:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nd := niceDate(&tt.tm)
			if nd != tt.want {
				t.Errorf("want %q; got %q", "17-Dec-20 04:00", nd)
			}
		})
	}

}
