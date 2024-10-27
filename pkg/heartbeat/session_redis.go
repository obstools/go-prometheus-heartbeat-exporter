package heartbeat

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type sessionRedis struct {
	*sessionHeartbeat
}

// Redis heartbeat session logic implementation

func (session *sessionRedis) run() error {
	var err error

	opts, err := redis.ParseURL(session.url + connectionRedisProtocol)
	if err != nil {
		return err
	}
	rdb := redis.NewClient(opts)
	defer rdb.Close()

	if session.query != "" {
		cmds, ctx := strings.Split(session.query, ";"), context.Background()
		for _, cmd := range cmds {
			parts := strings.Fields(strings.TrimSpace(cmd))
			if len(parts) < 2 {
				return fmt.Errorf("missed argument for redis command: %s", cmd)
			}

			switch strings.ToUpper(parts[0]) {
			case "SET":
				if len(parts) >= 3 {
					err = rdb.Set(ctx, parts[1], parts[2], 0).Err()
				}
			case "GET":
				_, err = rdb.Get(ctx, parts[1]).Result()
			case "DEL":
				err = rdb.Del(ctx, parts[1]).Err()
			default:
				err = fmt.Errorf("unknown redis command: %s", parts[0])
			}
		}
	} else {
		err = rdb.Ping(context.Background()).Err()
	}

	return err
}
