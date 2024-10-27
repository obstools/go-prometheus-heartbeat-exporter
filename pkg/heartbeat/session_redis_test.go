package heartbeat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionRedisRun(t *testing.T) {
	url := "redis://localhost:6379"
	t.Run("when wrong session url provided", func(t *testing.T) {
		assert.NotNil(t, createNewSession(connectionRedis, url+"/ololo", "").run())
	})
	t.Run("when connection is established", func(t *testing.T) {
		assert.Nil(t, createNewSession(connectionRedis, url, "SET key1 value1; GET key1; DEL key1").run())
	})

	t.Run("when ping fails", func(t *testing.T) {
		assert.NotNil(t, createNewSession(connectionRedis, "redis://user:password@localhost:6379", "").run())
	})

	t.Run("when unknown command provided", func(t *testing.T) {
		command := "OLOLO"
		err := createNewSession(connectionRedis, url, command+" 42 43").run()

		assert.EqualError(t, err, fmt.Sprintf("unknown redis command: %s", command))
	})

	t.Run("when broken query provided", func(t *testing.T) {
		assert.EqualError(t, createNewSession(connectionRedis, url, "SET").run(), "missed argument for redis command: SET")
	})
}
