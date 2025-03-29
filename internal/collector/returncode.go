package collector

import (
	"github.com/zinrai/prom-textfile-exporter/internal/config"
	"github.com/zinrai/prom-textfile-exporter/internal/executor"
)

// collects metrics based on command return codes
type ReturnCodeCollector struct {
	metricConfig config.MetricConfig
	timeoutSec   int
}

// creates a new ReturnCodeCollector
func NewReturnCodeCollector(metricConfig config.MetricConfig, timeoutSec int) *ReturnCodeCollector {
	return &ReturnCodeCollector{
		metricConfig: metricConfig,
		timeoutSec:   timeoutSec,
	}
}

// executes the command and returns its exit code as the metric value
func (c *ReturnCodeCollector) Collect() CollectResult {
	collector := c.metricConfig.Collector

	// Execute command
	result := executor.ExecuteCommandWithResult(
		collector.Command,
		c.timeoutSec,
	)

	// Create metric with the exit code, regardless of whether the command succeeded
	metric := Metric{
		Name:   c.metricConfig.Name,
		Value:  float64(result.ExitCode),
		Type:   c.metricConfig.Type,
		Help:   c.metricConfig.Help,
		Labels: collector.Labels,
	}

	// For returncode collectors, metrics are always valid regardless of errors
	return CollectResult{
		Metric:      metric,
		MetricValid: true,
		Error:       result.Error,
		HasWarning:  result.Error != nil,
	}
}
