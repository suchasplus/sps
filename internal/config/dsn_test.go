package config

import (
	"os"
	"testing"
)

func TestGetDSN(t *testing.T) {
	// Test case 1: DSN from file
	t.Run("from file", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test-dsn-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		expectedDSN := "file_dsn_string"
		if _, err := tmpfile.WriteString(expectedDSN); err != nil {
			t.Fatal(err)
		}
		tmpfile.Close()

		dsn, err := GetDSN(tmpfile.Name(), "arg_dsn_string")
		if err != nil {
			t.Errorf("GetDSN() error = %v, wantErr nil", err)
		}
		if dsn != expectedDSN {
			t.Errorf("GetDSN() got = %v, want %v", dsn, expectedDSN)
		}
	})

	// Test case 2: DSN from argument
	t.Run("from argument", func(t *testing.T) {
		expectedDSN := "arg_dsn_string"
		dsn, err := GetDSN("", expectedDSN)
		if err != nil {
			t.Errorf("GetDSN() error = %v, wantErr nil", err)
		}
		if dsn != expectedDSN {
			t.Errorf("GetDSN() got = %v, want %v", dsn, expectedDSN)
		}
	})

	// Test case 3: No DSN provided
	t.Run("no dsn provided", func(t *testing.T) {
		_, err := GetDSN("", "")
		if err == nil {
			t.Error("GetDSN() error = nil, wantErr not nil")
		}
	})

	// Test case 4: File not found
	t.Run("file not found", func(t *testing.T) {
		_, err := GetDSN("non_existent_file.txt", "")
		if err == nil {
			t.Error("GetDSN() error = nil, wantErr not nil")
		}
	})
}
