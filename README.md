# MasstDB

A powerful command-line database backup utility supporting multiple database management systems.

## Features

- **Multi-Database Support** - PostgreSQL, MySQL, MongoDB, SQLite
- **Compression** - Gzip compression enabled by default
- **Simple CLI** - Easy-to-use commands for backup and restore
- **Cross-Platform** - Works on macOS, Linux, and Windows

## Installation

### Using Go

```bash
go install github.com/adityanarayan/masstdb@latest
```

### From Source

```bash
git clone https://github.com/adityanarayan/masstdb.git
cd masstdb
make build
```

## Quick Start

### Backup a Database

```bash
# SQLite
masstdb backup --type sqlite --database ./mydata.db

# PostgreSQL
masstdb backup --type postgres --host localhost --user admin --password secret --database mydb

# MySQL
masstdb backup --type mysql --host localhost --user root --password secret --database mydb

# MongoDB
masstdb backup --type mongodb --host localhost --database mydb
```

### Restore a Database

```bash
# Restore from backup file
masstdb restore --type sqlite --database ./mydata.db --file backups/mydata.db_full_20260130.sql.gz

# Restore PostgreSQL
masstdb restore --type postgres --host localhost --user admin --password secret --database mydb --file backups/mydb_full_20260130.sql.gz
```

### List Backups

```bash
masstdb list

# Output:
# NAME                                 SIZE   CREATED
# ----                                 ----   -------
# mydb_full_20260130_152700.sql.gz     208 B  2026-01-30 15:27:00
```

### Test Connection

```bash
masstdb test --type postgres --host localhost --user admin --password secret --database mydb
```

## Command Reference

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | | Path to config file (default: `$HOME/.masstdb.yaml`) |
| `--verbose` | `-v` | Enable verbose output |
| `--help` | `-h` | Show help |

### Backup Command

```bash
masstdb backup [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--type` | `-t` | required | Database type (postgres, mysql, mongodb, sqlite) |
| `--host` | `-H` | localhost | Database host |
| `--port` | `-P` | auto | Database port |
| `--user` | `-u` | | Database username |
| `--password` | `-p` | | Database password |
| `--database` | `-d` | required | Database name or path |
| `--output` | `-o` | ./backups | Output directory |
| `--compress` | `-c` | true | Compress backup with gzip |
| `--backup-type` | `-b` | full | Backup type (full, incremental, differential) |

### Restore Command

```bash
masstdb restore [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--type` | `-t` | Database type (required) |
| `--file` | `-f` | Backup file to restore (required) |
| `--database` | `-d` | Target database (required) |
| `--tables` | | Specific tables to restore (comma-separated) |

### List Command

```bash
masstdb list [flags]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | `-d` | ./backups | Directory to list backups from |

## Examples

### Daily Backup Script

```bash
#!/bin/bash
# backup.sh - Run daily with cron

DATE=$(date +%Y%m%d)
masstdb backup --type postgres \
  --host localhost \
  --user backup_user \
  --password "$DB_PASSWORD" \
  --database production \
  --output /var/backups/db
```

### Backup Without Compression

```bash
masstdb backup --type sqlite --database ./app.db --compress=false
```

### Backup to Custom Directory

```bash
masstdb backup --type postgres --database mydb --output /mnt/backup-drive/databases
```

## Default Ports

| Database | Default Port |
|----------|-------------|
| PostgreSQL | 5432 |
| MySQL | 3306 |
| MongoDB | 27017 |
| SQLite | N/A (file-based) |

## Prerequisites

MasstDB uses native database tools. Ensure these are installed:

| Database | Required Tools |
|----------|---------------|
| PostgreSQL | `pg_dump`, `psql` |
| MySQL | `mysqldump`, `mysql` |
| MongoDB | `mongodump`, `mongorestore` |
| SQLite | `sqlite3` |

## Configuration File

Create `~/.masstdb.yaml` for default settings:

```yaml
default_database:
  type: postgres
  host: localhost
  username: admin

storage:
  local_path: ./backups

backup:
  compress: true
  default_type: full
```

## Building from Source

```bash
# Build for current platform
make build

# Build for all platforms
make release

# Run tests
make test

# Clean build artifacts
make clean
```

## Project Structure

```
masstdb/
├── main.go                 # Entry point
├── cmd/                    # CLI commands
│   ├── root.go            # Root command
│   ├── backup.go          # Backup command
│   ├── restore.go         # Restore command
│   ├── list.go            # List command
│   └── test_connection.go # Test command
├── internal/
│   ├── database/          # Database connectors
│   ├── backup/            # Backup service
│   ├── config/            # Configuration
│   └── logger/            # Logging
└── Makefile               # Build automation
```

## License

MIT License

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Author

Aditya Narayan
