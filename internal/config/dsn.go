package config

import (
	"errors"
	"os"
)

// GetDSN resolves the database source name (DSN) from the command-line arguments.
// It prioritizes the 'file' flag over the direct DSN argument.
func GetDSN(filePath, dsnArg string) (string, error) {
	// If filePath is provided, read DSN from the file.
	if filePath != "" {
		dsnBytes, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return string(dsnBytes), nil
	}

	// If no file is provided, use the direct DSN argument.
	if dsnArg != "" {
		return dsnArg, nil
	}

	// If neither is provided, return an error.
	return "", errors.New("no DSN provided; use the -f flag or provide a DSN argument")
}
