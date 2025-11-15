# atlas-migrate-status

A simple CLI tool to view complete [Atlas](https://atlasgo.io/) migration history with execution times.

## Why?

`atlas migrate status` only shows a summary:

```
Migration Status: OK
  -- Current Version: 20250923175718
  -- Next Version:    Already at latest version
  -- Executed Files:  5
  -- Pending Files:   0
```

But what if you want to see **all** migrations with their execution times? That's where `atlas-migrate-status` comes in:

```
Migration History (5 total)
────────────────────────────────────────────────────────────────────────────────────────────────────
VERSION         DESCRIPTION              EXECUTED AT          DURATION  TYPE      STATUS
20250920100000  Add users table         2025-09-20 10:05:32  45ms      baseline  ✓
20250921150000  Add indexes             2025-09-21 15:12:18  230ms     versioned ✓
20250922120000  Add foreign keys        2025-09-22 12:34:56  1.23s     versioned ✓
20250923175718  Create job_postings     2025-09-23 17:57:28  12ms      versioned ✓
```

## Installation

### Quick Install (Recommended)

**macOS / Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/okm321/atlas-migrate-status/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
# Download the latest release from:
# https://github.com/okm321/atlas-migrate-status/releases
```

### Using Go

```bash
go install github.com/okm321/atlas-migrate-status@latest
```

### Manual Download

Download the latest binary for your platform from [Releases](https://github.com/okm321/atlas-migrate-status/releases).

## Usage

### Basic Usage

```bash
atlas-migrate-status --url "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

### Using atlas.hcl Environment

If you have an `atlas.hcl` file with environment configurations:

```hcl
env "local" {
  url = "postgres://user:pass@localhost:5432/mydb"
}

env "prod" {
  url = "postgres://user:pass@prod-db:5432/mydb"
  migration {
    revisions_schema = "public.atlas_revisions"  # Optional: custom table name
  }
}
```

You can use the `--env` flag:

```bash
# Use the "local" environment
atlas-migrate-status --env local

# Use the "prod" environment
atlas-migrate-status --env prod

# Specify custom config file path
atlas-migrate-status --env local --config /path/to/atlas.hcl
```

See [examples/atlas.hcl](examples/atlas.hcl) for a complete example configuration.

### Custom Revisions Table

If you're using a custom schema or table name:

```bash
atlas-migrate-status --url "postgres://..." --revisions-schema my_schema.my_revisions
```

### Verbose Mode

```bash
atlas-migrate-status --url "postgres://..." --verbose
```

### Help

```bash
atlas-migrate-status --help
```

## Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--url` | `-u` | - | Database URL (required if --env not used) |
| `--env` | `-e` | - | Environment from atlas.hcl (required if --url not used) |
| `--config` | `-c` | `./atlas.hcl` | Path to atlas.hcl config file |
| `--revisions-schema` | - | `atlas_schema_revisions` | Revisions table name |
| `--verbose` | `-v` | `false` | Enable verbose output |

## How It Works

Atlas stores migration history in the `atlas_schema_revisions` table (by default). This tool queries that table and displays the data in a human-readable format.

The query is essentially:

```sql
SELECT 
  version,
  description,
  executed_at,
  execution_time,
  type,
  error
FROM atlas_schema_revisions
ORDER BY executed_at ASC;
```

## Requirements

- PostgreSQL database with Atlas migrations applied
- The `atlas_schema_revisions` table must exist (created by `atlas migrate apply`)

## Supported Databases

- ✅ PostgreSQL
- ⏳ MySQL (coming soon)
- ⏳ SQLite (coming soon)

## Development

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (for testing)

### Build

```bash
make build
```

### Install dependencies

```bash
go mod download
go mod tidy
```

### Run tests

```bash
make test
```

### Run locally

```bash
go run . --url "postgres://user:pass@localhost:5432/dbname"
```

### Format code

```bash
make fmt
```

### See all available commands

```bash
make help
```

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## Roadmap

- [x] Basic functionality
- [x] Custom revisions table support
- [x] Verbose mode
- [x] Support `--env` flag to read from `atlas.hcl`
- [ ] MySQL support
- [ ] SQLite support
- [ ] Filter options (`--last N`, `--after DATE`)

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [Atlas](https://atlasgo.io/) - The amazing database schema migration tool
- [tablewriter](https://github.com/olekukonko/tablewriter) - ASCII table rendering
- [Cobra](https://github.com/spf13/cobra) - CLI framework
