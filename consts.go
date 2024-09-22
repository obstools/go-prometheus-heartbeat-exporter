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

	// Metric
	metricLabelInstanceName       = "instance_name"
	metricLabelErrorType          = "error_type"
	metricNameSuccessfulAttempts  = "heartbeat_successful_attempts"
	metricNameFailedAttempts      = "heartbeat_failed_attempts"
	metricNameDuration            = "heartbeat_duration_sec"
	metricDescSuccessfulAttempts  = "The total number of successful heartbeat attempts"
	metricDescFailedAttempts      = "The total number of failed heartbeat attempts"
	metricDescDuration            = "The heartbeat duration in seconds"
	metricFailureLableConncection = "connection"
	metricFailureLableTimeout     = "timeout"

	// Session
	connectionWithoutSsl = "?sslmode=disable"
)
