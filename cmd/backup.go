package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AdityaNarayan29/masstDB/internal/backup"
	"github.com/AdityaNarayan29/masstDB/internal/database"
	"github.com/AdityaNarayan29/masstDB/internal/logger"
	"github.com/spf13/cobra"
)

var (
	// Database connection flags
	dbType   string
	host     string
	port     int
	username string
	password string
	dbName   string

	// Backup options
	outputDir  string
	compress   bool
	backupType string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a database backup",
	Long: `Create a backup of the specified database.

Supported database types: postgres, mysql, mongodb, sqlite

Backup types:
  - full: Complete backup of the entire database
  - incremental: Only changes since the last backup (if supported)
  - differential: Changes since the last full backup (if supported)

Examples:
  # Backup PostgreSQL database
  dbbackup backup --type postgres --host localhost --port 5432 --user admin --password secret --database mydb

  # Backup with compression
  dbbackup backup --type postgres --database mydb --compress

  # Backup SQLite database
  dbbackup backup --type sqlite --database /path/to/database.db`,
	RunE: runBackup,
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Database connection flags
	backupCmd.Flags().StringVarP(&dbType, "type", "t", "", "database type (postgres, mysql, mongodb, sqlite)")
	backupCmd.Flags().StringVarP(&host, "host", "H", "localhost", "database host")
	backupCmd.Flags().IntVarP(&port, "port", "P", 0, "database port (default depends on db type)")
	backupCmd.Flags().StringVarP(&username, "user", "u", "", "database username")
	backupCmd.Flags().StringVarP(&password, "password", "p", "", "database password")
	backupCmd.Flags().StringVarP(&dbName, "database", "d", "", "database name")

	// Backup options
	backupCmd.Flags().StringVarP(&outputDir, "output", "o", "./backups", "output directory for backup files")
	backupCmd.Flags().BoolVarP(&compress, "compress", "c", true, "compress backup file")
	backupCmd.Flags().StringVarP(&backupType, "backup-type", "b", "full", "backup type (full, incremental, differential)")

	// Mark required flags
	backupCmd.MarkFlagRequired("type")
	backupCmd.MarkFlagRequired("database")
}

func runBackup(cmd *cobra.Command, args []string) error {
	log := logger.New(verbose)
	log.Info("Starting backup process...")

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

	// Test connection
	log.Info("Testing database connection...")
	if err := connector.TestConnection(); err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	log.Info("Connection successful!")

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate backup filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s_%s", dbName, backupType, timestamp)
	outputPath := filepath.Join(outputDir, filename)

	// Create backup service
	backupService := backup.NewService(log)

	// Perform backup
	log.Info("Creating backup...")
	startTime := time.Now()

	result, err := backupService.Backup(connector, backup.Options{
		Type:       backupType,
		OutputPath: outputPath,
		Compress:   compress,
	})
	if err != nil {
		log.Error("Backup failed: %v", err)
		return fmt.Errorf("backup failed: %w", err)
	}

	duration := time.Since(startTime)

	// Log results
	log.Info("Backup completed successfully!")
	log.Info("  File: %s", result.FilePath)
	log.Info("  Size: %s", formatBytes(result.Size))
	log.Info("  Duration: %s", duration.Round(time.Millisecond))

	return nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
