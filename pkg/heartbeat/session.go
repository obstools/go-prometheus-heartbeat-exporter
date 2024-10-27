package heartbeat

// Session base heartbeat structure
type sessionHeartbeat struct {
	connection string
	url        string
	query      string
}

// Session type
type session interface {
	run() error
	getConnection() string
	getURL() string
	getQuery() string
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

// Returns session query
func (session *sessionHeartbeat) getQuery() string {
	return session.query
}

// New session builder. Returns new session interface based on the connection type
func newSession(connection, url, query string) session {
	sessionHeartbeat := &sessionHeartbeat{url: url, query: query}

	switch connection {
	case connectionPostgres:
		sessionHeartbeat.connection = connection
		return &sessionPostgres{sessionHeartbeat: sessionHeartbeat}
	case connectionRedis:
		sessionHeartbeat.connection = connection
		return &sessionRedis{sessionHeartbeat: sessionHeartbeat}
	}

	return nil
}
