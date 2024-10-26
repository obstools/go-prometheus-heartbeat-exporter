package heartbeat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("when configuration file is valid", func(t *testing.T) {
		server, err := New("../../fixtures/config_without_instances.yml")

		assert.NotNil(t, server)
		assert.NoError(t, err)
	})

	t.Run("when configuration file is invalid", func(t *testing.T) {
		server, err := New("../../fixtures/config_broken.yml")

		assert.Nil(t, server)
		assert.Error(t, err)
	})
}
