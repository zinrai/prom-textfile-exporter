package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zinrai/prom-textfile-exporter/internal/collector"
	"github.com/zinrai/prom-textfile-exporter/internal/config"
	"github.com/zinrai/prom-textfile-exporter/internal/writer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prom-textfile-exporter <command> [command options]")
		fmt.Println("\nCommands:")
		fmt.Println("  run       Execute metric collection")
		fmt.Println("  validate  Validate configuration file")
		fmt.Println("\nRun 'prom-textfile-exporter <command> -h' for help on a specific command")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		runCommand(os.Args[2:])
	case "validate":
		validateCommand(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: prom-textfile-exporter <command> [command options]")
		fmt.Println("\nCommands:")
		fmt.Println("  run       Execute metric collection")
		fmt.Println("  validate  Validate configuration file")
		os.Exit(1)
	}
}

func runCommand(args []string) {
	runFlags := flag.NewFlagSet("run", flag.ExitOnError)

	configFile := runFlags.String("config", "./config.yaml", "Path to configuration file")
	outputDir := runFlags.String("output-dir", "/tmp", "Output directory")
	timeoutSec := runFlags.Int("timeout", 10, "Command execution timeout in seconds")

	runFlags.Usage = func() {
		fmt.Println("Usage: prom-textfile-exporter run [options]")
		fmt.Println("\nOptions:")
		runFlags.PrintDefaults()
	}

	if err := runFlags.Parse(args); err != nil {
		fmt.Println(err)
		runFlags.Usage()
		os.Exit(1)
	}

	runExecute(*configFile, *outputDir, *timeoutSec)
}

func validateCommand(args []string) {
	validateFlags := flag.NewFlagSet("validate", flag.ExitOnError)

	configFile := validateFlags.String("config", "./config.yaml", "Path to configuration file")

	validateFlags.Usage = func() {
		fmt.Println("Usage: prom-textfile-exporter validate [options]")
		fmt.Println("\nOptions:")
		validateFlags.PrintDefaults()
	}

	if err := validateFlags.Parse(args); err != nil {
		fmt.Println(err)
		validateFlags.Usage()
		os.Exit(1)
	}

	validateExecute(*configFile)
}

func runExecute(configFile, outputDir string, timeoutSec int) {
	log.Printf("Loading configuration from: %s", configFile)

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	metrics := []collector.Metric{}
	hasErrors := false
	hasWarnings := false

	for name, metricCfg := range cfg.Metrics {
		log.Printf("Collecting metric: %s", name)

		col, err := collector.NewCollector(metricCfg, timeoutSec)
		if err != nil {
			log.Printf("Error creating collector for %s: %v", name, err)
			hasErrors = true
			continue
		}

		result := col.Collect()

		if result.Error != nil {
			// If there is an error but valid metrics
			if result.MetricValid {
				log.Printf("Warning collecting metric %s: %v", name, result.Error)
				hasWarnings = true
				metrics = append(metrics, result.Metric)
			} else {
				// If there is an error and no valid metrics
				log.Printf("Error collecting metric %s: %v", name, result.Error)
				hasErrors = true
			}
		} else if result.MetricValid {
			// エラーなしで有効なメトリクス
			metrics = append(metrics, result.Metric)
		}
	}

	if len(metrics) == 0 {
		log.Fatalf("No metrics were collected")
	}

	outputFile := filepath.Join(outputDir, "prom_textfile_exporter.prom")
	err = writer.WriteMetrics(metrics, outputFile)
	if err != nil {
		log.Fatalf("Failed to write metrics to file: %v", err)
	}

	log.Printf("Successfully wrote %d metrics to %s", len(metrics), outputFile)
	if hasWarnings {
		log.Printf("Some warnings occurred during collection, but metrics were still generated")
	}
	if hasErrors {
		log.Printf("Some errors occurred during collection, not all metrics were generated")
	}
}

func validateExecute(configFile string) {
	log.Printf("Validating configuration in: %s", configFile)

	_, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	fmt.Println("Configuration is valid.")
}
