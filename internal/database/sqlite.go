package database

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// SQLiteConnector implements database operations for SQLite
type SQLiteConnector struct {
	config Config
}

// NewSQLiteConnector creates a new SQLite connector
func NewSQLiteConnector(config Config) (*SQLiteConnector, error) {
	return &SQLiteConnector{config: config}, nil
}

// TestConnection tests if the SQLite database file exists and is accessible
func (s *SQLiteConnector) TestConnection() error {
	// Check if file exists
	info, err := os.Stat(s.config.Database)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("database file does not exist: %s", s.config.Database)
		}
		return fmt.Errorf("cannot access database file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a database file: %s", s.config.Database)
	}

	// Try to open the database with sqlite3
	cmd := exec.Command("sqlite3", s.config.Database, "SELECT 1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to open database: %s - %s", err, string(output))
	}

	return nil
}

// Backup performs a SQLite backup using .dump command
func (s *SQLiteConnector) Backup(w io.Writer) error {
	// Use sqlite3 .dump command to create SQL backup
	cmd := exec.Command("sqlite3", s.config.Database, ".dump")
	cmd.Stdout = w

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sqlite3 dump failed: %s - %s", err, stderr.String())
	}

	return nil
}

// Restore restores a SQLite database from backup
func (s *SQLiteConnector) Restore(r io.Reader) error {
	// Use sqlite3 to execute the SQL dump
	cmd := exec.Command("sqlite3", s.config.Database)
	cmd.Stdin = r

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restore failed: %s - %s", err, stderr.String())
	}

	return nil
}

// Close closes the SQLite connection
func (s *SQLiteConnector) Close() error {
	return nil
}

// Type returns the database type
func (s *SQLiteConnector) Type() string {
	return "sqlite"
}

// SupportsIncremental returns true if incremental backups are supported
func (s *SQLiteConnector) SupportsIncremental() bool {
	return false
}
