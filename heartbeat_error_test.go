package heartbeat

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeartbeatErrorOriginalError(t *testing.T) {
	t.Run("returns original error", func(t *testing.T) {
		err := errors.New("test error")
		heartbeatError := heartbeatError{err: err}

		assert.Equal(t, err, heartbeatError.originalError())
	})
}
