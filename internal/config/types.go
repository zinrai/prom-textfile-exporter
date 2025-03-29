package config

type Config struct {
	Metrics map[string]MetricConfig `yaml:"metrics"`
}

type MetricConfig struct {
	Name      string          `yaml:"name"`
	Type      string          `yaml:"type"`
	Help      string          `yaml:"help"`
	Collector CollectorConfig `yaml:"collector"`
}

type CollectorConfig struct {
	Type    string             `yaml:"type"`
	Command string             `yaml:"command"`
	Labels  map[string]string  `yaml:"labels"`
	Mapping map[string]float64 `yaml:"mapping,omitempty"`
	Parse   *ParseConfig       `yaml:"parse,omitempty"`
}

type ParseConfig struct {
	Pattern      string             `yaml:"pattern"`
	Index        int                `yaml:"index"`
	ValueType    string             `yaml:"value_type,omitempty"`    // "float", "int", "bool", "string" etc...
	StringMap    map[string]float64 `yaml:"string_map,omitempty"`    // String-to-number mapping
	Multiplier   float64            `yaml:"multiplier,omitempty"`    // Numeric value to multiply the extracted value by
	DefaultValue *float64           `yaml:"default_value,omitempty"` // Default value if parsing fails
}
