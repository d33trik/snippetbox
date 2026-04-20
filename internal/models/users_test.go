package models

import (
	"testing"

	"codeberg.org/d33trik/snippetbox/internal/assert"
)

func TestUserModelExistsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := map[string]struct {
		userID int
		want   bool
	}{
		"Valid ID": {
			userID: 1,
			want:   true,
		},
		"Zero ID": {
			userID: 0,
			want:   false,
		},
		"Non-existent ID": {
			userID: 2,
			want:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tc.userID)
			assert.Equal(t, exists, tc.want)
			assert.Nil(t, err)
		})
	}
}
