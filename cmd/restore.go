package cmd

import (
	"fmt"
	"time"

	"github.com/AdityaNarayan29/masstDB/internal/backup"
	"github.com/AdityaNarayan29/masstDB/internal/database"
	"github.com/AdityaNarayan29/masstDB/internal/logger"
	"github.com/spf13/cobra"
)

var (
	// Restore specific flags
	backupFile string
	tables     []string
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a database from backup",
	Long: `Restore a database from a backup file.

The restore command supports:
  - Full database restoration
  - Selective table restoration (if supported by DBMS)
  - Automatic decompression of .gz files

Examples:
  # Restore full database
  dbbackup restore --file backup.sql.gz --type postgres --database mydb

  # Restore specific tables
  dbbackup restore --file backup.sql.gz --type postgres --database mydb --tables users,orders`,
	RunE: runRestore,
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	// Database connection flags (reusing from backup)
	restoreCmd.Flags().StringVarP(&dbType, "type", "t", "", "database type (postgres, mysql, mongodb, sqlite)")
	restoreCmd.Flags().StringVarP(&host, "host", "H", "localhost", "database host")
	restoreCmd.Flags().IntVarP(&port, "port", "P", 0, "database port")
	restoreCmd.Flags().StringVarP(&username, "user", "u", "", "database username")
	restoreCmd.Flags().StringVarP(&password, "password", "p", "", "database password")
	restoreCmd.Flags().StringVarP(&dbName, "database", "d", "", "database name")

	// Restore specific flags
	restoreCmd.Flags().StringVarP(&backupFile, "file", "f", "", "backup file to restore from")
	restoreCmd.Flags().StringSliceVar(&tables, "tables", nil, "specific tables to restore (comma-separated)")

	// Mark required flags
	restoreCmd.MarkFlagRequired("type")
	restoreCmd.MarkFlagRequired("database")
	restoreCmd.MarkFlagRequired("file")
}

func runRestore(cmd *cobra.Command, args []string) error {
	log := logger.New(verbose)
	log.Info("Starting restore process...")

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

	// Create database connector
	connector, err := database.NewConnector(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to create database connector: %w", err)
	}

	// Test connection (skip for SQLite as file may not exist yet)
	if dbType != "sqlite" {
		log.Info("Testing database connection...")
		if err := connector.TestConnection(); err != nil {
			return fmt.Errorf("connection test failed: %w", err)
		}
		log.Info("Connection successful!")
	} else {
		log.Info("SQLite restore - will create database file if needed")
	}

	// Create backup service
	backupService := backup.NewService(log)

	// Perform restore
	log.Info("Restoring from: %s", backupFile)
	startTime := time.Now()

	err = backupService.Restore(connector, backup.RestoreOptions{
		FilePath: backupFile,
		Tables:   tables,
	})
	if err != nil {
		log.Error("Restore failed: %v", err)
		return fmt.Errorf("restore failed: %w", err)
	}

	duration := time.Since(startTime)

	log.Info("Restore completed successfully!")
	log.Info("  Duration: %s", duration.Round(time.Millisecond))

	return nil
}
