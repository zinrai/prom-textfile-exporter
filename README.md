# prom-textfile-exporter

A tool for generating Prometheus metrics for the Node Exporter textfile collector.

## Overview

`prom-textfile-exporter` allows you to collect metrics from various commands and their outputs using a YAML configuration file. It supports several collection methods:

1. **Command return codes** - Use exit codes directly as metric values
2. **Return code mapping** - Map exit codes to specific metric values
3. **Output parsing** - Extract values from command output using regular expressions

## Installation

```bash
$ go build -o prom-textfile-exporter cmd/main.go
```

## Usage

### Basic Usage

```bash
$ prom-textfile-exporter run --config /path/to/config.yaml
```

### Command-line Options

```
prom-textfile-exporter <command> [command options]

Commands:
  run       Execute metric collection
  validate  Validate configuration file

Command Options (run):
  --config <path>       Path to configuration file (default: ./config.yaml)
  --output-dir <path>   Output directory (default: /tmp)
  --timeout <seconds>   Command execution timeout in seconds (default: 10)

Command Options (validate):
  --config <path>       Path to configuration file (default: ./config.yaml)
```

## Configuration

`prom-textfile-exporter` uses YAML files for configuration. See the examples configuration file.

## Integration with Node Exporter

To use with the Node Exporter's textfile collector:

1. Set up a cron job to run `prom-textfile-exporter` periodically
2. Configure Node Exporter with `--collector.textfile.directory=/path/to/metrics/dir`

Example cron job:

```
*/5 * * * * /usr/local/bin/prom-textfile-exporter run --config /etc/prom-textfile-exporter/config.yaml
```

## Error Handling

Metrics are generated even if the command fails.

- For the `returncode` collector, the exit code is always captured as the metric value
- For the `returncode_mapping` collector, the exit code is mapped to a value based on configuration
- For the `output_parse` collector, a default value can be specified for use when parsing fails

## License

This project is licensed under the [MIT License](./LICENSE).
