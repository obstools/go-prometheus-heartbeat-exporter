package heartbeat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	t.Run("returns new postgres session interface", func(t *testing.T) {
		connection, url := "postgres", "some_url"
		session := newSession(connection, url)

		assert.NotNil(t, session)
		assert.Equal(t, connection, session.getConnection())
		assert.Equal(t, url, session.getURL())
	})

	t.Run("returns nil for undefined connection", func(t *testing.T) {
		assert.Nil(t, newSession("undefined", "some_url"))
	})
}
