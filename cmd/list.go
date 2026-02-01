package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var (
	listDir string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available backups",
	Long: `List all backup files in the specified directory.

Examples:
  dbbackup list
  dbbackup list --dir /path/to/backups`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listDir, "dir", "d", "./backups", "directory to list backups from")
}

type backupInfo struct {
	Name    string
	Size    int64
	ModTime time.Time
}

func runList(cmd *cobra.Command, args []string) error {
	// Check if directory exists
	if _, err := os.Stat(listDir); os.IsNotExist(err) {
		fmt.Printf("No backups found. Directory '%s' does not exist.\n", listDir)
		return nil
	}

	// Read directory
	entries, err := os.ReadDir(listDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Filter and collect backup files
	var backups []backupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if it's a backup file (by extension)
		name := entry.Name()
		if !isBackupFile(name) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, backupInfo{
			Name:    name,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	if len(backups) == 0 {
		fmt.Printf("No backups found in '%s'\n", listDir)
		return nil
	}

	// Sort by modification time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].ModTime.After(backups[j].ModTime)
	})

	// Print table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSIZE\tCREATED")
	fmt.Fprintln(w, "----\t----\t-------")

	for _, b := range backups {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			b.Name,
			formatBytes(b.Size),
			b.ModTime.Format("2006-01-02 15:04:05"),
		)
	}

	w.Flush()

	fmt.Printf("\nTotal: %d backup(s) in %s\n", len(backups), listDir)
	fullPath, _ := filepath.Abs(listDir)
	fmt.Printf("Location: %s\n", fullPath)

	return nil
}

func isBackupFile(name string) bool {
	extensions := []string{".sql", ".sql.gz", ".dump", ".dump.gz", ".bson", ".bson.gz", ".db", ".db.gz"}
	for _, ext := range extensions {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}
	return false
}
