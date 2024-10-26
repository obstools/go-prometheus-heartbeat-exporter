package heartbeat

import (
	"database/sql"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/lib/pq"
)

var testDB = "heartbeat_test"

// Regex builder
func newRegex(regexPattern string) (*regexp.Regexp, error) {
	return regexp.Compile(regexPattern)
}

// Returns log message regex based on log level and message context
func loggerMessageRegex(logLevel, logMessage string) *regexp.Regexp {
	regex, _ := newRegex(logLevel + `: \d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}\.\d{6} \[` + logMessage + "]")
	return regex
}

// Postgres session
// Creates test database for postgres
func createPostgresDb() error {
	db, err := sql.Open("postgres", "host=localhost user=postgres port=5432 sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + testDB)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P04" {
			return nil
		}
	}
	return err
}

// Drops test database for postgres
func dropPostgresDb() error {
	db, err := sql.Open("postgres", "host=localhost user=postgres port=5432 sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS " + testDB)
	return err
}

// Composes connection string for postgres
func composePostgresConnectionString() string {
	return "postgres://postgres@localhost:5432/" + testDB
}

// Creates new session
func createNewSession(connection, url, query string) session {
	return newSession(connection, url, query)
}

// Creates new wait group
func createNewWaitGroup() *sync.WaitGroup {
	return new(sync.WaitGroup)
}

// Creates new minimal configuration
func createNewMinimalConfiguration() *Configuration {
	return &Configuration{
		Port:         8080,
		MetricsRoute: "/metrics",
	}
}

// Generates a unique instance name for testing purposes
func generateUniqueInstanceName() string {
	return fmt.Sprintf("test_instance_%d", time.Now().UnixNano())
}

// New heartbeat instance
func newTestInstance(session session) *heartbeatInstance {
	return &heartbeatInstance{
		name:        generateUniqueInstanceName(),
		intervalSec: 2,
		timeoutSec:  3,
		session:     session,
	}
}
