package heartbeat

// Heartbeat error wrapper structure
type heartbeatError struct {
	duration float64
	err      error
}

// heartbeatError methods

// Returns original (unwrapped) error
func (errorWrapper *heartbeatError) originalError() error {
	return errorWrapper.err
}
