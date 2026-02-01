package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Used for flags
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "masstdb",
	Short: "A database backup utility supporting multiple DBMS",
	Long: `MasstDB is a CLI tool for backing up and restoring databases.

Supported databases:
  - PostgreSQL
  - MySQL
  - MongoDB
  - SQLite

Features:
  - Full, incremental, and differential backups
  - Compression support (gzip)
  - Local and cloud storage (AWS S3, GCS, Azure)
  - Backup scheduling
  - Detailed logging

Example usage:
  masstdb backup --type postgres --host localhost --database mydb
  masstdb restore --file backup_20240101.sql.gz --type postgres
  masstdb list --storage local`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.masstdb.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
