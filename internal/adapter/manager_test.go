package adapter

import "testing"

func TestDetectDriver(t *testing.T) {
	testCases := []struct {
		name    string
		dsn     string
		want    string
		wantErr bool
	}{
		{"PostgreSQL", "postgres://user:pass@host:5432/db", "postgres", false},
		{"PostgreSQL Alternative", "postgresql://user:pass@host:5432/db", "postgres", false},
		{"MySQL", "user:pass@tcp(host:3306)/db", "mysql", false},
		{"SQLite", "/path/to/my.db", "sqlite3", false},
		{"SQLite Alternative", "file.sqlite", "sqlite3", false},
		{"SQLite 3", "data.sqlite3", "sqlite3", false},
		{"Unsupported", "mongodb://user:pass@host/db", "", true},
		{"Empty DSN", "", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DetectDriver(tc.dsn)
			if (err != nil) != tc.wantErr {
				t.Errorf("DetectDriver() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("DetectDriver() = %v, want %v", got, tc.want)
			}
		})
	}
}
