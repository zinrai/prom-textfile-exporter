package collector

import (
	"fmt"
	"regexp"

	"github.com/zinrai/prom-textfile-exporter/internal/config"
	"github.com/zinrai/prom-textfile-exporter/internal/executor"
)

// collects metrics by parsing command output
type OutputParseCollector struct {
	metricConfig config.MetricConfig
	timeoutSec   int
}

// creates a new OutputParseCollector
func NewOutputParseCollector(metricConfig config.MetricConfig, timeoutSec int) (*OutputParseCollector, error) {
	// Validate parse configuration
	if metricConfig.Collector.Parse == nil {
		return nil, fmt.Errorf("output_parse collector requires parse configuration")
	}

	return &OutputParseCollector{
		metricConfig: metricConfig,
		timeoutSec:   timeoutSec,
	}, nil
}

// executes the command and parses its output for the metric value
func (c *OutputParseCollector) Collect() CollectResult {
	collector := c.metricConfig.Collector
	parse := collector.Parse

	result := CollectResult{
		MetricValid: false,
	}

	metric := Metric{
		Name:   c.metricConfig.Name,
		Type:   c.metricConfig.Type,
		Help:   c.metricConfig.Help,
		Labels: collector.Labels,
	}

	cmdResult := executor.ExecuteCommandWithResult(
		collector.Command,
		c.timeoutSec,
	)

	if cmdResult.Error != nil {
		// If default values are set
		if parse.DefaultValue != nil {
			metric.Value = *parse.DefaultValue
			result.Metric = metric
			result.MetricValid = true
			result.HasWarning = true
			result.Error = fmt.Errorf("command execution failed (using default value): %w", cmdResult.Error)
			return result
		}

		// If there is no default value
		result.Error = fmt.Errorf("command execution failed: %w", cmdResult.Error)
		return result
	}

	output := cmdResult.Output
	if output == "" {
		if parse.DefaultValue != nil {
			metric.Value = *parse.DefaultValue
			result.Metric = metric
			result.MetricValid = true
			result.HasWarning = true
			result.Error = fmt.Errorf("empty command output (using default value)")
			return result
		}

		result.Error = fmt.Errorf("empty command output")
		return result
	}

	re, err := regexp.Compile(parse.Pattern)
	if err != nil {
		if parse.DefaultValue != nil {
			metric.Value = *parse.DefaultValue
			result.Metric = metric
			result.MetricValid = true
			result.HasWarning = true
			result.Error = fmt.Errorf("invalid regex pattern (using default value): %w", err)
			return result
		}

		result.Error = fmt.Errorf("invalid regex pattern: %w", err)
		return result
	}

	matches := re.FindStringSubmatch(output)
	if len(matches) <= parse.Index {
		if parse.DefaultValue != nil {
			metric.Value = *parse.DefaultValue
			result.Metric = metric
			result.MetricValid = true
			result.HasWarning = true
			result.Error = fmt.Errorf("pattern didn't match or index out of range (using default value)")
			return result
		}

		result.Error = fmt.Errorf("pattern didn't match or index out of range")
		return result
	}

	// Value Extraction and Conversion
	extractedStr := matches[parse.Index]
	value, err := convertValue(extractedStr, parse)
	if err != nil {
		if parse.DefaultValue != nil {
			metric.Value = *parse.DefaultValue
			result.Metric = metric
			result.MetricValid = true
			result.HasWarning = true
			result.Error = fmt.Errorf("could not parse value (using default): %w", err)
			return result
		}

		result.Error = fmt.Errorf("could not parse value from output: %w", err)
		return result
	}

	// If the value is successfully obtained
	metric.Value = value
	result.Metric = metric
	result.MetricValid = true

	return result
}
