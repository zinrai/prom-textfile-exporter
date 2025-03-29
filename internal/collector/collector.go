package collector

import (
	"fmt"

	"github.com/zinrai/prom-textfile-exporter/internal/config"
)

type Metric struct {
	Name   string
	Value  float64
	Type   string
	Help   string
	Labels map[string]string
}

type CollectResult struct {
	Metric      Metric // Metrics collected
	MetricValid bool   // Metrics valid
	Error       error  // Errors encountered
	HasWarning  bool   // Are there any warnings, e.g., if default values are used
}

type Collector interface {
	// Collect executes the command and returns the result containing metric and error info
	Collect() CollectResult
}

func NewCollector(metricConfig config.MetricConfig, timeoutSec int) (Collector, error) {
	switch metricConfig.Collector.Type {
	case "returncode":
		return NewReturnCodeCollector(metricConfig, timeoutSec), nil
	case "returncode_mapping":
		return NewReturnCodeMappingCollector(metricConfig, timeoutSec)
	case "output_parse":
		return NewOutputParseCollector(metricConfig, timeoutSec)
	default:
		return nil, fmt.Errorf("unknown collector type: %s", metricConfig.Collector.Type)
	}
}
