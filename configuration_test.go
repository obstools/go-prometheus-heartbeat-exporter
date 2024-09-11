package heartbeat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	t.Run("when configuration file exists and valid", func(t *testing.T) {
		expectedConfiguration := &Configuration{
			LogToStdout:  true,
			LogActivity:  true,
			Port:         8080,
			MetricsRoute: "/metrics",
			InstancesAttributes: []*InstanceAttributes{
				{
					Name:        "postgres_1",
					Connection:  "postgres",
					URL:         "postgres://localhost:5432/healthcheck_db",
					IntervalSec: 3,
					TimeoutSec:  2,
				},
			},
		}
		actualConfiguration, err := loadConfiguration("./fixtures/config.yml")

		assert.EqualValues(t, expectedConfiguration, actualConfiguration)
		assert.Nil(t, err)
	})

	t.Run("when configuration file not exists", func(t *testing.T) {
		configurationPath := "./fixtures/not_existing_config.yml"
		actualConfiguration, err := loadConfiguration(configurationPath)

		assert.EqualError(t, err, "open "+configurationPath+": no such file or directory")
		assert.Nil(t, actualConfiguration)
	})

	t.Run("when configuration file is not a file", func(t *testing.T) {
		configurationPath := "./fixtures"
		actualConfiguration, err := loadConfiguration(configurationPath)

		assert.EqualError(t, err, "yaml: input error: read "+configurationPath+": is a directory")
		assert.Nil(t, actualConfiguration)
	})
}
