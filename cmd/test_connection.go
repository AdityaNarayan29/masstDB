package cmd

import (
	"fmt"

	"github.com/AdityaNarayan29/masstDB/internal/database"
	"github.com/AdityaNarayan29/masstDB/internal/logger"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test database connection",
	Long: `Test the connection to a database without performing any backup.

This is useful for validating connection parameters before running a backup.

Examples:
  dbbackup test --type postgres --host localhost --user admin --password secret --database mydb`,
	RunE: runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Database connection flags
	testCmd.Flags().StringVarP(&dbType, "type", "t", "", "database type (postgres, mysql, mongodb, sqlite)")
	testCmd.Flags().StringVarP(&host, "host", "H", "localhost", "database host")
	testCmd.Flags().IntVarP(&port, "port", "P", 0, "database port")
	testCmd.Flags().StringVarP(&username, "user", "u", "", "database username")
	testCmd.Flags().StringVarP(&password, "password", "p", "", "database password")
	testCmd.Flags().StringVarP(&dbName, "database", "d", "", "database name")

	// Mark required flags
	testCmd.MarkFlagRequired("type")
	testCmd.MarkFlagRequired("database")
}

func runTest(cmd *cobra.Command, args []string) error {
	log := logger.New(verbose)

	// Set default ports based on database type
	if port == 0 {
		port = database.DefaultPort(dbType)
	}

	// Create database configuration
	dbConfig := database.Config{
		Type:     dbType,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: dbName,
	}

	// Validate configuration
	if err := dbConfig.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log.Info("Testing connection to %s database '%s' at %s:%d...",
		dbConfig.Type, dbConfig.Database, dbConfig.Host, dbConfig.Port)

	// Create database connector
	connector, err := database.NewConnector(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to create database connector: %w", err)
	}

	// Test connection
	if err := connector.TestConnection(); err != nil {
		log.Error("Connection failed: %v", err)
		return fmt.Errorf("connection test failed: %w", err)
	}

	log.Info("Connection successful!")
	return nil
}
