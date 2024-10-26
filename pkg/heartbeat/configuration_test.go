package heartbeat

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpolateEnvVars(t *testing.T) {
	defer os.Unsetenv("TEST_ENV")

	t.Run("when ENV var found in the string", func(t *testing.T) {
		targetValue := "some_text"
		os.Setenv("TEST_ENV", targetValue)

		assert.Equal(t, "abc"+targetValue+"def", interpolateEnvVars("abc${TEST_ENV}def"))
	})
}

func TestLoadConfiguration(t *testing.T) {
	defer os.Unsetenv("METRICS_ROUTE")

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
					URL:         "postgres://localhost:5432/heartbeat_test",
					Query:       "CREATE TABLE tmp (id SERIAL PRIMARY KEY); DROP TABLE tmp",
					IntervalSec: 3,
					TimeoutSec:  2,
				},
			},
		}
		actualConfiguration, err := loadConfiguration("../../fixtures/config.yml")

		assert.EqualValues(t, expectedConfiguration, actualConfiguration)
		assert.Nil(t, err)
	})

	t.Run("when configuration file exists and valid", func(t *testing.T) {
		additionalRoute := "/some_path"
		os.Setenv("METRICS_ROUTE", additionalRoute)

		expectedConfiguration := &Configuration{
			LogToStdout:  true,
			LogActivity:  true,
			Port:         8080,
			MetricsRoute: "/metrics" + additionalRoute,
			InstancesAttributes: []*InstanceAttributes{
				{
					Name:        "postgres_1",
					Connection:  "postgres",
					URL:         "postgres://localhost:5432/heartbeat_test",
					IntervalSec: 3,
					TimeoutSec:  2,
				},
			},
		}
		actualConfiguration, err := loadConfiguration("../../fixtures/config_with_env.yml")

		assert.EqualValues(t, expectedConfiguration, actualConfiguration)
		assert.Nil(t, err)
	})

	t.Run("when configuration file not exists", func(t *testing.T) {
		configurationPath := "../../fixtures/config_not_existing.yml"
		actualConfiguration, err := loadConfiguration(configurationPath)

		assert.EqualError(t, err, "open "+configurationPath+": no such file or directory")
		assert.Nil(t, actualConfiguration)
	})

	t.Run("when configuration file is a broken file", func(t *testing.T) {
		configurationPath := "../../fixtures/config"
		actualConfiguration, err := loadConfiguration(configurationPath)

		assert.EqualError(t, err, "yaml: invalid trailing UTF-8 octet")
		assert.Nil(t, actualConfiguration)
	})

	t.Run("when configuration file is not a file", func(t *testing.T) {
		configurationPath := "../../fixtures"
		actualConfiguration, err := loadConfiguration(configurationPath)

		assert.EqualError(t, err, "read "+configurationPath+": is a directory")
		assert.Nil(t, actualConfiguration)
	})

	t.Run("when configuration file is broken", func(t *testing.T) {
		configurationPath := "../../fixtures/config_broken.yml"
		actualConfiguration, err := loadConfiguration(configurationPath)

		assert.EqualError(t, err, "yaml: line 7: did not find expected key")
		assert.Nil(t, actualConfiguration)
	})
}
