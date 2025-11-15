package cmd

import (
	"fmt"
	"os"

	"github.com/okm321/atlas-migrate-status/internal/config"
	"github.com/okm321/atlas-migrate-status/internal/db"
	"github.com/okm321/atlas-migrate-status/internal/display"
	"github.com/spf13/cobra"
)

var (
	url             string
	env             string
	configPath      string
	revisionsSchema string
	verbose         bool
	Version         = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "atlas-migrate-status",
	Short: "View all Atlas migration history",
	Long: `atlas-migrate-status displays complete migration history from Atlas schema revisions table.

Unlike 'atlas migrate status' which shows only a summary, this tool displays
all migrations with their execution times, timestamps, and status.

Examples:
  # Show migration history using database URL
  atlas-migrate-status --url "postgres://user:pass@localhost:5432/dbname"

  # Use environment from atlas.hcl
  atlas-migrate-status --env local

  # Use specific config file
  atlas-migrate-status --env local --config /path/to/atlas.hcl

  # Output as JSON
  atlas-migrate-status --url "postgres://..." --format json`,
	RunE:    runCommand,
	Version: Version,
}

func runCommand(cmd *cobra.Command, args []string) error {
	if url == "" && env == "" {
		return fmt.Errorf("either --url or --env must be specified\n\nUsage: atlas-migrate-status --url <database-url>")
	}

	if url != "" && env != "" {
		return fmt.Errorf("--url and --env are mutually exclusive, use only one")
	}

	dbURL := url
	if env != "" {
		// Load config from atlas.hcl
		if verbose {
			if configPath != "" {
				fmt.Fprintf(os.Stderr, "Loading config from: %s\n", configPath)
			} else {
				fmt.Fprintf(os.Stderr, "Looking for atlas.hcl...\n")
			}
		}

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w\n\nMake sure atlas.hcl exists in the current directory or use --config to specify the path", err)
		}

		envConfig, err := cfg.GetEnv(env)
		if err != nil {
			return fmt.Errorf("failed to get environment config: %w", err)
		}

		if envConfig.URL == "" {
			return fmt.Errorf("no URL configured for environment '%s'", env)
		}

		dbURL = envConfig.URL

		if envConfig.RevisionsSchema != "" && !cmd.Flags().Changed("revisions-schema") {
			revisionsSchema = envConfig.RevisionsSchema
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Using environment: %s\n", env)
			fmt.Fprintf(os.Stderr, "Database URL: %s\n", maskPassword(dbURL))
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Connecting to database...\n")
		fmt.Fprintf(os.Stderr, "Revisions table: %s\n", revisionsSchema)
	}

	migrations, err := db.FetchMigrationHistory(dbURL, revisionsSchema)
	if err != nil {
		return fmt.Errorf("failed to fetch migration history: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Found %d migrations\n\n", len(migrations))
	}

	display.PrintTable(migrations)

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// maskPassword masks the password in a database URL for display
func maskPassword(dbURL string) string {
	// Simple masking: postgres://user:pass@host -> postgres://user:****@host
	// This is just for display purposes in verbose mode
	var masked string
	inPassword := false
	for i, c := range dbURL {
		if c == ':' && i > 0 && dbURL[i-1] != '/' {
			inPassword = true
			masked += string(c)
			continue
		}
		if inPassword && c == '@' {
			inPassword = false
			masked += "****" + string(c)
			continue
		}
		if !inPassword {
			masked += string(c)
		}
	}
	return masked
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "Database URL (postgres://user:pass@localhost:5432/dbname)")
	rootCmd.Flags().StringVarP(&env, "env", "e", "", "Environment from atlas.hcl")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to atlas.hcl (default: ./atlas.hcl)")
	rootCmd.Flags().StringVar(&revisionsSchema, "revisions-schema", "atlas_schema_revisions", "Schema revisions table name")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
