package heartbeat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	t.Run("returns new metric", func(t *testing.T) {
		metric := newMetric("test_instance")

		assert.NotNil(t, metric)
	})
}
