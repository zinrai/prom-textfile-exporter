package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// loads and validates the configuration from a file
func LoadConfig(filename string) (*Config, error) {
	// Read configuration file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validates the configuration
func validateConfig(config *Config) error {
	if len(config.Metrics) == 0 {
		return fmt.Errorf("no metrics defined")
	}

	// Validate each metric
	for name, metric := range config.Metrics {
		if err := validateMetric(name, metric); err != nil {
			return fmt.Errorf("invalid metric '%s': %w", name, err)
		}
	}

	return nil
}

// validates a single metric configuration
func validateMetric(name string, metric MetricConfig) error {
	// Validate metric name
	if metric.Name == "" {
		return fmt.Errorf("metric name is required")
	}

	// Validate metric type
	if metric.Type == "" {
		return fmt.Errorf("metric type is required")
	}
	if metric.Type != "gauge" && metric.Type != "counter" {
		return fmt.Errorf("metric type must be 'gauge' or 'counter', got '%s'", metric.Type)
	}

	// Validate collector
	return validateCollector(metric.Collector)
}

// validates the collector configuration
func validateCollector(collector CollectorConfig) error {
	// Check command
	if collector.Command == "" {
		return fmt.Errorf("collector command is required")
	}

	// Check collector type
	switch collector.Type {
	case "returncode":
		// No additional validation needed
	case "returncode_mapping":
		if collector.Mapping == nil || len(collector.Mapping) == 0 {
			return fmt.Errorf("returncode_mapping requires a mapping configuration")
		}
	case "output_parse":
		if err := validateParseConfig(collector.Parse); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown collector type: %s", collector.Type)
	}

	return nil
}

// validates the parse configuration
func validateParseConfig(parse *ParseConfig) error {
	if parse == nil {
		return fmt.Errorf("parse configuration is required")
	}
	if parse.Pattern == "" {
		return fmt.Errorf("parse pattern is required")
	}
	// Check if pattern is a valid regular expression
	_, err := regexp.Compile(parse.Pattern)
	if err != nil {
		return fmt.Errorf("invalid regular expression pattern: %w", err)
	}
	if parse.Index < 0 {
		return fmt.Errorf("parse index must be non-negative")
	}
	return nil
}
