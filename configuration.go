package heartbeat

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Instance attributes structure
type InstanceAttributes struct {
	Name        string `yaml:"name"`
	Connection  string `yaml:"connection"`
	URL         string `yaml:"url"`
	IntervalSec int    `yaml:"interval"`
	TimeoutSec  int    `yaml:"timeout"`
}

// Configuration structure
type Configuration struct {
	LogToStdout         bool                  `yaml:"log_to_stdout"`
	LogActivity         bool                  `yaml:"log_activity"`
	Port                int                   `yaml:"port"`
	MetricsRoute        string                `yaml:"metrics_route"`
	ShutdownTimeout     int                   `yaml:"shutdown_timeout"`
	InstancesAttributes []*InstanceAttributes `yaml:"instances"`
}

// Configuration loader/builder. Returns configuration pointer and error
func loadConfiguration(configurationPath string) (configuration *Configuration, err error) {
	file, err := os.Open(configurationPath)
	if err != nil {
		return configuration, err
	}
	defer file.Close()

	if file != nil {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&configuration); err != nil {
			return configuration, err
		}
	}

	// TODO: add defaults configuration values in the next release
	// TODO: add configuration validation (unique instance names) in the next release

	return
}
