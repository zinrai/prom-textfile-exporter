package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zinrai/prom-textfile-exporter/internal/collector"
)

// writes metrics in Prometheus format to a file
func WriteMetrics(metrics []collector.Metric, outputFile string) error {
	// Build string builder for the content
	var sb strings.Builder

	// Process each metric
	uniqueMetrics := make(map[string]bool)

	for _, metric := range metrics {
		// Add HELP and TYPE lines only once per metric name
		if !uniqueMetrics[metric.Name] {
			fmt.Fprintf(&sb, "# HELP %s %s\n", metric.Name, metric.Help)
			fmt.Fprintf(&sb, "# TYPE %s %s\n", metric.Name, metric.Type)
			uniqueMetrics[metric.Name] = true
		}

		// Format labels if any
		labelsStr := ""
		if len(metric.Labels) > 0 {
			var labelParts []string
			for k, v := range metric.Labels {
				labelParts = append(labelParts, fmt.Sprintf("%s=%q", k, v))
			}
			labelsStr = fmt.Sprintf("{%s}", strings.Join(labelParts, ","))
		}

		// Add metric line
		fmt.Fprintf(&sb, "%s%s %g\n", metric.Name, labelsStr, metric.Value)
	}

	// Create a temporary file
	dir := filepath.Dir(outputFile)
	tmpfile, err := os.CreateTemp(dir, "metrics.*.prom")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	// Clean up on error
	defer func() {
		tmpfile.Close()
		if err != nil {
			os.Remove(tmpfile.Name())
		}
	}()

	// Write content to temporary file
	if _, err := tmpfile.WriteString(sb.String()); err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// Close the file before renaming
	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Atomically move the file to its final destination
	if err := os.Rename(tmpfile.Name(), outputFile); err != nil {
		return fmt.Errorf("failed to move temporary file: %w", err)
	}

	return nil
}
