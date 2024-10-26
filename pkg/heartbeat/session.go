package heartbeat

// Session base heartbeat structure
type sessionHeartbeat struct {
	connection string
	url        string
}

// Session type
type session interface {
	run() error
	getConnection() string
	getURL() string
}

// Session methods

// Returns session connection name
func (session *sessionHeartbeat) getConnection() string {
	return session.connection
}

// Returns session url
func (session *sessionHeartbeat) getURL() string {
	return session.url
}

// New session builder. Returns new session interface based on the connection type
func newSession(connection, url string) session {
	switch connection {
	case connectionPostgres:
		return &sessionPostgres{
			sessionHeartbeat: &sessionHeartbeat{
				connection: connection,
				url:        url,
			},
		}
	}

	return nil
}
