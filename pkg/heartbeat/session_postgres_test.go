package heartbeat

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSessionPostgresRun(t *testing.T) {
	t.Run("returns nil if connection is established", func(t *testing.T) {
		if err := createPostgresDb(); err != nil {
			t.Fatalf("Failed to create test database: %v", err)
		}
		defer func() {
			if err := dropPostgresDb(); err != nil {
				t.Fatalf("Failed to drop test database: %v", err)
			}
		}()

		assert.Nil(t, createNewSession("postgres", composePostgresConnectionString()).run())
	})

	t.Run("returns error if ping fails", func(t *testing.T) {
		assert.NotNil(t, createNewSession("postgres", "postgres://user:password@localhost:5432").run())
	})

	t.Run("returns error if connection is not established", func(t *testing.T) {
		session := &sessionPostgres{
			sessionHeartbeat: &sessionHeartbeat{
				connection: "ologres",
				url:        composePostgresConnectionString(),
			},
		}

		assert.NotNil(t, session.run())
	})
}
