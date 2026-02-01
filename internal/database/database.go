package database

import (
	"fmt"
	"io"
)

// Config holds database connection configuration
type Config struct {
	Type     string // postgres, mysql, mongodb, sqlite
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// Validate checks if the configuration is valid
func (c Config) Validate() error {
	if c.Type == "" {
		return fmt.Errorf("database type is required")
	}

	validTypes := map[string]bool{
		"postgres": true,
		"mysql":    true,
		"mongodb":  true,
		"sqlite":   true,
	}

	if !validTypes[c.Type] {
		return fmt.Errorf("unsupported database type: %s", c.Type)
	}

	if c.Database == "" {
		return fmt.Errorf("database name is required")
	}

	// SQLite doesn't need host/port/credentials
	if c.Type == "sqlite" {
		return nil
	}

	if c.Host == "" {
		return fmt.Errorf("host is required for %s", c.Type)
	}

	return nil
}

// ConnectionString returns a formatted connection string
func (c Config) ConnectionString() string {
	switch c.Type {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.Username, c.Password, c.Database)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			c.Username, c.Password, c.Host, c.Port, c.Database)
	case "mongodb":
		if c.Username != "" {
			return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
				c.Username, c.Password, c.Host, c.Port, c.Database)
		}
		return fmt.Sprintf("mongodb://%s:%d/%s", c.Host, c.Port, c.Database)
	case "sqlite":
		return c.Database // SQLite uses file path directly
	default:
		return ""
	}
}

// DefaultPort returns the default port for a database type
func DefaultPort(dbType string) int {
	ports := map[string]int{
		"postgres": 5432,
		"mysql":    3306,
		"mongodb":  27017,
		"sqlite":   0, // SQLite doesn't use a port
	}
	return ports[dbType]
}

// Connector defines the interface for database operations
type Connector interface {
	// TestConnection tests if the database connection works
	TestConnection() error

	// Backup performs a database backup and writes to the provided writer
	Backup(w io.Writer) error

	// Restore restores a database from the provided reader
	Restore(r io.Reader) error

	// Close closes any open connections
	Close() error

	// Type returns the database type
	Type() string

	// SupportsIncremental returns true if incremental backups are supported
	SupportsIncremental() bool
}

// NewConnector creates a new database connector based on the configuration
func NewConnector(config Config) (Connector, error) {
	switch config.Type {
	case "postgres":
		return NewPostgresConnector(config)
	case "mysql":
		return NewMySQLConnector(config)
	case "sqlite":
		return NewSQLiteConnector(config)
	case "mongodb":
		return NewMongoDBConnector(config)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}
