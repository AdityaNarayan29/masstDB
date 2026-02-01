package database

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// PostgresConnector implements database operations for PostgreSQL
type PostgresConnector struct {
	config Config
}

// NewPostgresConnector creates a new PostgreSQL connector
func NewPostgresConnector(config Config) (*PostgresConnector, error) {
	return &PostgresConnector{config: config}, nil
}

// TestConnection tests the PostgreSQL connection
func (p *PostgresConnector) TestConnection() error {
	// Build psql command for connection test
	args := p.buildPsqlArgs()
	args = append(args, "-c", "SELECT 1")

	cmd := exec.Command("psql", args...)
	cmd.Env = p.buildEnv()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("connection failed: %s - %s", err, string(output))
	}

	return nil
}

// Backup performs a PostgreSQL backup using pg_dump
func (p *PostgresConnector) Backup(w io.Writer) error {
	// Build pg_dump command
	args := []string{
		"-h", p.config.Host,
		"-p", fmt.Sprintf("%d", p.config.Port),
		"-U", p.config.Username,
		"-d", p.config.Database,
		"-F", "p", // plain text format
		"--no-password",
	}

	cmd := exec.Command("pg_dump", args...)
	cmd.Stdout = w
	cmd.Env = p.buildEnv()

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %s - %s", err, stderr.String())
	}

	return nil
}

// Restore restores a PostgreSQL database from backup
func (p *PostgresConnector) Restore(r io.Reader) error {
	// Build psql command for restore
	args := p.buildPsqlArgs()

	cmd := exec.Command("psql", args...)
	cmd.Stdin = r
	cmd.Env = p.buildEnv()

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restore failed: %s - %s", err, stderr.String())
	}

	return nil
}

// Close closes the PostgreSQL connection (no persistent connection to close)
func (p *PostgresConnector) Close() error {
	return nil
}

// Type returns the database type
func (p *PostgresConnector) Type() string {
	return "postgres"
}

// SupportsIncremental returns true if incremental backups are supported
func (p *PostgresConnector) SupportsIncremental() bool {
	return false // Standard pg_dump doesn't support incremental backups
}

// buildPsqlArgs builds common psql command arguments
func (p *PostgresConnector) buildPsqlArgs() []string {
	return []string{
		"-h", p.config.Host,
		"-p", fmt.Sprintf("%d", p.config.Port),
		"-U", p.config.Username,
		"-d", p.config.Database,
		"--no-password",
	}
}

// buildEnv builds environment variables for pg commands
func (p *PostgresConnector) buildEnv() []string {
	return []string{
		fmt.Sprintf("PGPASSWORD=%s", p.config.Password),
	}
}
