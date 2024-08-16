package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    time.Time
		expectError bool
	}{
		{
			name:        "Valid time string",
			input:       "2024-07-07T06:41:19.168000000Z",
			expected:    time.Date(2024, 7, 7, 6, 41, 19, 168000000, time.UTC),
			expectError: false,
		},
		{
			name:        "Invalid time string format",
			input:       "invalid time format",
			expected:    time.Time{},
			expectError: true,
		},
		{
			name:        "Time with year less than 1970",
			input:       "1969-12-31T23:59:59Z",
			expected:    time.Time{},
			expectError: true,
		},
		{
			name:        "Time with exact 1970",
			input:       "1970-01-01T00:00:00Z",
			expected:    time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedTime, err := ParseTime(tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, parsedTime)
			}
		})
	}
}
