package heartbeat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	url, query := "some_url", "some_query"

	t.Run("when connection is postgres", func(t *testing.T) {
		connection := "postgres"
		session := newSession(connection, url, query)

		assert.NotNil(t, session)
		assert.Equal(t, connection, session.getConnection())
		assert.Equal(t, url, session.getURL())
		assert.Equal(t, query, session.getQuery())
	})

	t.Run("when connection is redis", func(t *testing.T) {
		connection := "redis"
		session := newSession(connection, url, query)

		assert.NotNil(t, session)
		assert.Equal(t, connection, session.getConnection())
		assert.Equal(t, url, session.getURL())
		assert.Equal(t, query, session.getQuery())
	})

	t.Run("returns nil for undefined connection", func(t *testing.T) {
		assert.Nil(t, newSession("undefined", url, query))
	})
}

func TestSessionGetConnection(t *testing.T) {
	connection := "some_connection"
	session := &sessionHeartbeat{connection: connection}

	assert.Equal(t, connection, session.getConnection())
}

func TestSessionGetURL(t *testing.T) {
	url := "some_url"
	session := &sessionHeartbeat{url: url}

	assert.Equal(t, url, session.getURL())
}

func TestSessionGetQuery(t *testing.T) {
	query := "some_query"
	session := &sessionHeartbeat{query: query}

	assert.Equal(t, query, session.getQuery())
}
