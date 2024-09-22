package heartbeat

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSessionPostgres(t *testing.T) {
	t.Run("returns nil if connection is established", func(t *testing.T) {
		if err := createPostgresDb(); err != nil {
			t.Fatalf("Failed to create test database: %v", err)
		}
		defer func() {
			if err := dropPostgresDb(); err != nil {
				t.Fatalf("Failed to drop test database: %v", err)
			}
		}()

		assert.Nil(t, sessionPostgres(composePostgresConnectionString()))
	})

	t.Run("returns error if connection is not established", func(t *testing.T) {
		connectionPostgres = "olologres"
		defer func() {
			connectionPostgres = "postgres"
		}()
		err := sessionPostgres("postgres://localhost:5432")

		assert.NotNil(t, err)
	})
	t.Run("returns error if ping fails", func(t *testing.T) {
		err := sessionPostgres("postgres://user:password@localhost:5432")

		assert.NotNil(t, err)
	})
}
