package infrastructure

import (
	"testing"
)

// This is a bit tricky to test because the error message is not very useful
func TestColorMsg(t *testing.T) {
	testCases := []struct {
		name     string
		msg      string
		color    string
		expected string
	}{
		{
			name:     "cyan",
			msg:      "test",
			color:    cyan,
			expected: "\033[1;36mtest\033[0m",
		},
		{
			name:     "magenta",
			msg:      "test",
			color:    magenta,
			expected: "\033[1;35mtest\033[0m",
		},
		{
			name:     "blue",
			msg:      "test",
			color:    blue,
			expected: "\033[1;34mtest\033[0m",
		},
		{
			name:     "yellow",
			msg:      "test",
			color:    yellow,
			expected: "\033[1;33mtest\033[0m",
		},
		{
			name:     "red",
			msg:      "test",
			color:    red,
			expected: "\033[1;31mtest\033[0m",
		},
		{
			name:     "white",
			msg:      "test",
			color:    white,
			expected: "\033[1;37mtest\033[0m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// When
			result := ColorMsg(tc.msg, tc.color)

			// THEN
			if result != tc.expected {
				t.Errorf("Expected: %s\nGot: %s", tc.expected, result)
			}
		})
	}
}
