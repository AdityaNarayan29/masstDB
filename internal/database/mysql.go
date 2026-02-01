package database

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// MySQLConnector implements database operations for MySQL
type MySQLConnector struct {
	config Config
}

// NewMySQLConnector creates a new MySQL connector
func NewMySQLConnector(config Config) (*MySQLConnector, error) {
	return &MySQLConnector{config: config}, nil
}

// TestConnection tests the MySQL connection
func (m *MySQLConnector) TestConnection() error {
	args := m.buildMysqlArgs()
	args = append(args, "-e", "SELECT 1")

	cmd := exec.Command("mysql", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("connection failed: %s - %s", err, string(output))
	}

	return nil
}

// Backup performs a MySQL backup using mysqldump
func (m *MySQLConnector) Backup(w io.Writer) error {
	args := []string{
		"-h", m.config.Host,
		"-P", fmt.Sprintf("%d", m.config.Port),
		"-u", m.config.Username,
		fmt.Sprintf("-p%s", m.config.Password),
		"--single-transaction", // Consistent backup without locking
		"--routines",           // Include stored procedures
		"--triggers",           // Include triggers
		m.config.Database,
	}

	cmd := exec.Command("mysqldump", args...)
	cmd.Stdout = w

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysqldump failed: %s - %s", err, stderr.String())
	}

	return nil
}

// Restore restores a MySQL database from backup
func (m *MySQLConnector) Restore(r io.Reader) error {
	args := m.buildMysqlArgs()

	cmd := exec.Command("mysql", args...)
	cmd.Stdin = r

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restore failed: %s - %s", err, stderr.String())
	}

	return nil
}

// Close closes the MySQL connection
func (m *MySQLConnector) Close() error {
	return nil
}

// Type returns the database type
func (m *MySQLConnector) Type() string {
	return "mysql"
}

// SupportsIncremental returns true if incremental backups are supported
func (m *MySQLConnector) SupportsIncremental() bool {
	return false // Standard mysqldump doesn't support incremental backups
}

// buildMysqlArgs builds common mysql command arguments
func (m *MySQLConnector) buildMysqlArgs() []string {
	return []string{
		"-h", m.config.Host,
		"-P", fmt.Sprintf("%d", m.config.Port),
		"-u", m.config.Username,
		fmt.Sprintf("-p%s", m.config.Password),
		m.config.Database,
	}
}
