package backup

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AdityaNarayan29/masstDB/internal/database"
	"github.com/AdityaNarayan29/masstDB/internal/logger"
)

// Options contains backup configuration options
type Options struct {
	Type       string // full, incremental, differential
	OutputPath string
	Compress   bool
}

// RestoreOptions contains restore configuration options
type RestoreOptions struct {
	FilePath string
	Tables   []string // For selective restore
}

// Result contains information about a completed backup
type Result struct {
	FilePath string
	Size     int64
}

// Service handles backup and restore operations
type Service struct {
	log *logger.Logger
}

// NewService creates a new backup service
func NewService(log *logger.Logger) *Service {
	return &Service{log: log}
}

// Backup performs a database backup
func (s *Service) Backup(connector database.Connector, opts Options) (*Result, error) {
	// Determine output filename
	outputPath := opts.OutputPath
	extension := s.getExtension(connector.Type())
	outputPath += extension

	if opts.Compress {
		outputPath += ".gz"
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	var writer io.Writer = file

	// Add compression if requested
	var gzWriter *gzip.Writer
	if opts.Compress {
		gzWriter = gzip.NewWriter(file)
		defer gzWriter.Close()
		writer = gzWriter
	}

	// Perform backup
	s.log.Debug("Writing backup to: %s", outputPath)
	if err := connector.Backup(writer); err != nil {
		// Clean up failed backup file
		os.Remove(outputPath)
		return nil, err
	}

	// Ensure gzip is flushed
	if gzWriter != nil {
		if err := gzWriter.Close(); err != nil {
			return nil, fmt.Errorf("failed to close gzip writer: %w", err)
		}
	}

	// Get file info for size
	info, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &Result{
		FilePath: outputPath,
		Size:     info.Size(),
	}, nil
}

// Restore restores a database from backup
func (s *Service) Restore(connector database.Connector, opts RestoreOptions) error {
	// Open backup file
	file, err := os.Open(opts.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	var reader io.Reader = file

	// Check if file is compressed
	if strings.HasSuffix(opts.FilePath, ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	// Perform restore
	s.log.Debug("Restoring from: %s", opts.FilePath)
	if err := connector.Restore(reader); err != nil {
		return err
	}

	return nil
}

// getExtension returns the appropriate file extension for a database type
func (s *Service) getExtension(dbType string) string {
	switch dbType {
	case "postgres", "mysql", "sqlite":
		return ".sql"
	case "mongodb":
		return ".archive"
	default:
		return ".backup"
	}
}
