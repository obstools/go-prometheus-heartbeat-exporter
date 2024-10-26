package heartbeat

import (
	"io"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Instance attributes structure
type InstanceAttributes struct {
	Name        string `yaml:"name"`
	Connection  string `yaml:"connection"`
	URL         string `yaml:"url"`
	Query       string `yaml:"query"`
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

// Configuration loader helper. Returns new string with founded
// interpolated ENV vars matched with pattern ${VAR_NAME}
func interpolateEnvVars(source string) string {
	regex := regexp.MustCompile(`\$\{(\w+)\}`)

	return regex.ReplaceAllStringFunc(source, func(str string) string {
		varName := regex.FindStringSubmatch(str)[1]

		return os.Getenv(varName)
	})
}

// Configuration loader/builder. Returns configuration pointer and error
func loadConfiguration(configurationPath string) (configuration *Configuration, err error) {
	file, err := os.Open(configurationPath)
	if err != nil {
		return configuration, err
	}
	defer file.Close()

	configurationData, err := io.ReadAll(file)
	if err != nil {
		return configuration, err
	}

	err = yaml.Unmarshal([]byte(interpolateEnvVars(string(configurationData))), &configuration)
	if err != nil {
		return configuration, err
	}

	// TODO: add defaults configuration values in the next release
	// TODO: add configuration validation (unique instance names) in the next release

	return
}
