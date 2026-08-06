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

func TestGetMongoDBConfig(t *testing.T) {
	tests := []struct {
		name         string
		negativeTest bool
		uri          string
		dbName       string
	}{
		{
			name:         "invalid uri",
			negativeTest: true,
		},
		{
			name:         "valid uri, default db name",
			negativeTest: false,
			uri:          "http://localhost:27017",
		},
		{
			name:         "valid uri, valid db name",
			negativeTest: false,
			uri:          "http://localhost:27017",
			dbName:       "nextLayer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MONGODB_URI", tt.uri)
			if tt.dbName != "" {
				t.Setenv("MONGODB_DB_NAME", tt.dbName)
			}

			mongoDBConfig, err := GetMongoDBConfig()

			if tt.negativeTest {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if mongoDBConfig.URI != tt.uri {
				t.Errorf("expect %v, got %v", tt.uri, mongoDBConfig.URI)
			}

			if tt.dbName == "" {
				if mongoDBConfig.DBName != "baselayer" {
					t.Errorf("expect baselayer, got %v", mongoDBConfig.DBName)
				}
			} else {
				if mongoDBConfig.DBName != tt.dbName {
					t.Errorf("expect %v, got %v", tt.dbName, mongoDBConfig.DBName)
				}
			}
		})
	}
}
