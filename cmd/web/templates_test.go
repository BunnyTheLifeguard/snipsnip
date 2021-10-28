package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Berlin")

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 02, 02, 02, 0, 0, 0, time.UTC),
			want: "02 Feb 20 02:00 UTC",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2020, 02, 02, 02, 0, 0, 0, time.UTC).In(loc),
			want: "02 Feb 20 03:00 CET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			if hd != tt.want {
				t.Errorf("want %q got %q", tt.want, hd)
			}
		})
	}

}
