package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := map[string]struct {
		input time.Time
		want  string
	}{
		"Empty": {
			input: time.Time{},
			want:  "",
		},
		"UTC": {
			input: time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want:  "17 Mar 2024 at 10:15",
		},
		"CET": {
			input: time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want:  "17 Mar 2024 at 09:15",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := humanDate(tc.input)

			if got != tc.want {
				t.Errorf("got %q; want %q", got, tc.want)
			}
		})
	}
}
