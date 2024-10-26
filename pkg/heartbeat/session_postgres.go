package heartbeat

import "database/sql"

type sessionPostgres struct {
	*sessionHeartbeat
}

// Postgres heartbeat session logic implementation

func (session *sessionPostgres) run() error {
	var err error
	// time.Sleep(time.Duration(10) * time.Second) // Let's simulate that we exceeded expected timeout
	db, err := sql.Open(session.connection, session.url+connectionWithoutSsl)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	return err
}
