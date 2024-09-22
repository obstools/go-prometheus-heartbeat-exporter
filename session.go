package heartbeat

import "database/sql"

var connectionPostgres = "postgres"

// // Session type
// type session func(string) error

// Place for heartbeats logic implementation

// Postgres session
func sessionPostgres(url string) error {
	var err error
	// time.Sleep(time.Duration(10) * time.Second) // Let's simulate that we exceeded expected timeout
	db, err := sql.Open(connectionPostgres, url+connectionWithoutSsl)
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
