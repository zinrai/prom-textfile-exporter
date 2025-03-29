package collector

import (
	"fmt"
	"strconv"

	"github.com/zinrai/prom-textfile-exporter/internal/config"
	"github.com/zinrai/prom-textfile-exporter/internal/executor"
)

// collects metrics by mapping command return codes to values
type ReturnCodeMappingCollector struct {
	metricConfig config.MetricConfig
	mapping      map[int]float64
	defaultValue float64
	timeoutSec   int
}

// creates a new ReturnCodeMappingCollector
func NewReturnCodeMappingCollector(metricConfig config.MetricConfig, timeoutSec int) (*ReturnCodeMappingCollector, error) {
	// Convert string keys to int keys for easier lookup
	mapping := make(map[int]float64)
	var defaultValue float64
	hasDefault := false

	for k, v := range metricConfig.Collector.Mapping {
		if k == "default" {
			defaultValue = v
			hasDefault = true
			continue
		}

		exitCode, err := strconv.Atoi(k)
		if err != nil {
			return nil, fmt.Errorf("invalid exit code in mapping: %s", k)
		}
		mapping[exitCode] = v
	}

	if !hasDefault {
		return nil, fmt.Errorf("returncode_mapping requires a default value")
	}

	return &ReturnCodeMappingCollector{
		metricConfig: metricConfig,
		mapping:      mapping,
		defaultValue: defaultValue,
		timeoutSec:   timeoutSec,
	}, nil
}

// executes the command and maps its exit code to the metric value
func (c *ReturnCodeMappingCollector) Collect() CollectResult {
	collector := c.metricConfig.Collector

	// Execute command
	result := executor.ExecuteCommandWithResult(
		collector.Command,
		c.timeoutSec,
	)

	// Map exit code to value - always proceed regardless of error
	exitCode := result.ExitCode
	value, exists := c.mapping[exitCode]
	if !exists {
		value = c.defaultValue
	}

	// Create metric
	metric := Metric{
		Name:   c.metricConfig.Name,
		Value:  value,
		Type:   c.metricConfig.Type,
		Help:   c.metricConfig.Help,
		Labels: collector.Labels,
	}

	// For collectors, metrics are always valid regardless of errors
	return CollectResult{
		Metric:      metric,
		MetricValid: true,
		Error:       result.Error,
		HasWarning:  result.Error != nil,
	}
}
