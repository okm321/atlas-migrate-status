package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/okm321/atlas-migrate-status/internal/db"
	"github.com/olekukonko/tablewriter"
)

func PrintTable(migrations []db.Migration) {
	if len(migrations) == 0 {
		fmt.Println("No migrations found.")
		return
	}

	// Print summary header
	fmt.Printf("\nMigration History (%d total)\n", len(migrations))
	fmt.Println(strings.Repeat("─", 100))

	var data [][]string
	for _, m := range migrations {
		status := "✅"
		if m.Error != "" {
			status = "❌"
		}

		executedAt := m.ExecutedAt.Format("2006-01-02 15:04:05.000000")
		duration := formatDuration(m.ExecutionTime)

		description := m.Description
		if len(description) > 40 {
			description = description[:37] + "..."
		}

		data = append(data, []string{
			m.Version,
			description,
			executedAt,
			duration,
			m.Type,
			status,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Version", "Description", "Executed At", "Durataion", "Type", "Status"})
	table.Bulk(data)
	table.Render()
	fmt.Println()
}

func formatDuration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	seconds := float64(ms) / 1000.0
	return fmt.Sprintf("%.2fs", seconds)
}
