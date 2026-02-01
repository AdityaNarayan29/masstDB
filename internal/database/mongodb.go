package database

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// MongoDBConnector implements database operations for MongoDB
type MongoDBConnector struct {
	config Config
}

// NewMongoDBConnector creates a new MongoDB connector
func NewMongoDBConnector(config Config) (*MongoDBConnector, error) {
	return &MongoDBConnector{config: config}, nil
}

// TestConnection tests the MongoDB connection
func (m *MongoDBConnector) TestConnection() error {
	args := []string{
		m.config.ConnectionString(),
		"--eval", "db.runCommand({ ping: 1 })",
	}

	cmd := exec.Command("mongosh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try with legacy mongo shell
		cmd = exec.Command("mongo", args...)
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("connection failed: %s - %s", err, string(output))
		}
	}

	return nil
}

// Backup performs a MongoDB backup using mongodump
func (m *MongoDBConnector) Backup(w io.Writer) error {
	// mongodump writes to archive which we'll stream to the writer
	args := []string{
		"--host", m.config.Host,
		"--port", fmt.Sprintf("%d", m.config.Port),
		"--db", m.config.Database,
		"--archive", // Output to stdout as archive
	}

	if m.config.Username != "" {
		args = append(args, "--username", m.config.Username)
		args = append(args, "--password", m.config.Password)
		args = append(args, "--authenticationDatabase", "admin")
	}

	cmd := exec.Command("mongodump", args...)
	cmd.Stdout = w

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mongodump failed: %s - %s", err, stderr.String())
	}

	return nil
}

// Restore restores a MongoDB database from backup
func (m *MongoDBConnector) Restore(r io.Reader) error {
	args := []string{
		"--host", m.config.Host,
		"--port", fmt.Sprintf("%d", m.config.Port),
		"--db", m.config.Database,
		"--archive", // Read from stdin as archive
	}

	if m.config.Username != "" {
		args = append(args, "--username", m.config.Username)
		args = append(args, "--password", m.config.Password)
		args = append(args, "--authenticationDatabase", "admin")
	}

	cmd := exec.Command("mongorestore", args...)
	cmd.Stdin = r

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mongorestore failed: %s - %s", err, stderr.String())
	}

	return nil
}

// Close closes the MongoDB connection
func (m *MongoDBConnector) Close() error {
	return nil
}

// Type returns the database type
func (m *MongoDBConnector) Type() string {
	return "mongodb"
}

// SupportsIncremental returns true if incremental backups are supported
func (m *MongoDBConnector) SupportsIncremental() bool {
	return false // mongodump doesn't support incremental backups natively
}
