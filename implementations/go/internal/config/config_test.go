package config

import (
	"strings"
	"testing"
)

func TestGetPort(t *testing.T) {
	tests := []struct {
		name         string
		negativeTest bool
		port         string
		expected     string
	}{
		{
			name:     "valid port",
			port:     "8000",
			expected: "8000",
		},
		{
			name:         "conversion error",
			negativeTest: true,
			port:         "abc",
			expected:     "failed to parse port",
		},
		{
			name:         "range error",
			negativeTest: true,
			port:         "0",
			expected:     "invalid port",
		},
		{
			name:         "range error",
			negativeTest: true,
			port:         "65536",
			expected:     "invalid port",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.port)

			port, err := GetPort()

			if tt.negativeTest {
				if err == nil {
					t.Fatalf("expected an error, got nil")
				}

				if !strings.Contains(err.Error(), tt.expected) {
					t.Fatalf("expected an error containing %q, got %q", tt.expected, err.Error())
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if port != tt.expected {
				t.Errorf("expected port %q, got %q", tt.expected, port)
			}
		})
	}
}
