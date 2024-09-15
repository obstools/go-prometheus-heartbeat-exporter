package heartbeat

import "log"

const (
	// Exporter
	exporterStartMessage    = "Prometheus exporter started on "
	exporterErrorMsg        = "Failed to start prometheus exporter on port "
	exporterShutdownMessage = "Prometheus exporter is in the shutdown mode and won't accept new connections"
	exporterStopMessage     = "Prometheus exporter stopped gracefully"

	// Logger
	infoLogLevel    = "INFO"
	warningLogLevel = "WARNING"
	errorLogLevel   = "ERROR"
	logFlag         = log.Ldate | log.Lmicroseconds
)
