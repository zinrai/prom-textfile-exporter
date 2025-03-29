package collector

import (
	"fmt"
	"strconv"

	"github.com/zinrai/prom-textfile-exporter/internal/config"
)

// converts an extracted string to a float64 value based on parse configuration
func convertValue(str string, parse *config.ParseConfig) (float64, error) {
	var value float64
	var err error

	// If StringMap is defined, mapping is preferred
	if parse.StringMap != nil && len(parse.StringMap) > 0 {
		if mappedValue, ok := parse.StringMap[str]; ok {
			value = mappedValue
		} else if parse.DefaultValue != nil {
			value = *parse.DefaultValue
		} else {
			return 0, fmt.Errorf("string '%s' not found in mapping", str)
		}
	} else {
		switch parse.ValueType {
		case "", "float":
			value, err = strconv.ParseFloat(str, 64)
			if err != nil {
				return 0, fmt.Errorf("could not parse float value: %w", err)
			}
		case "int":
			intVal, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("could not parse int value: %w", err)
			}
			value = float64(intVal)
		case "bool":
			boolVal, err := strconv.ParseBool(str)
			if err != nil {
				return 0, fmt.Errorf("could not parse bool value: %w", err)
			}
			if boolVal {
				value = 1
			} else {
				value = 0
			}
		case "bool_nonzero":
			intVal, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("could not parse int value for bool_nonzero: %w", err)
			}
			if intVal != 0 {
				value = 1
			} else {
				value = 0
			}
		default:
			return 0, fmt.Errorf("unsupported value type: %s", parse.ValueType)
		}
	}

	// Multiplier applied ( 0 or not multiplied if not specified )
	if parse.Multiplier != 0 {
		value *= parse.Multiplier
	}

	return value, nil
}
