package heartbeat

import "log"

const (
	// Server
	serverStartMessage             = "Starting heartbeat workers, prometheus exporter..."
	serverStartErrorMessage        = "Unable to start heartbeat. Server must be inactive"
	serverStopMessage              = "Heartbeat workers, prometheus exporter were stopped gracefully"
	serverStopErrorMessage         = "Unable to stop heartbeat. Server must be active"
	serverStopExporterErrorMessage = "Errors during stopping prometheus exporter:"

	// Exporter
	exporterStartMessage    = "Prometheus exporter started on "
	exporterErrorMessage    = "Failed to start prometheus exporter on port "
	exporterShutdownMessage = "Prometheus exporter is in the shutdown mode and won't accept new connections"
	exporterStopMessage     = "Prometheus exporter stopped gracefully"
	exporterStatusTitle     = "Prometheus exporter status:"

	// Heartbeat
	heartbeatWorkerMsg                   = "Heartbeat worker:"
	heartbeatWorkerStatusActive          = "started"
	heartbeatWorkerStatusInactive        = "stopped"
	heartbeatActivitySuccess             = "SUCCESS:"
	heartbeatActivitySuccessElapsedTime  = "elapsed time:"
	heartbeatActivityFailure             = "FAILURE:"
	heartbeatActivityFailureSessionError = "heartbeat session error:"

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
	connectionPostgres      = "postgres"
	connectionWithoutSsl    = "?sslmode=disable"
	connectionRedis         = "redis"
	connectionRedisProtocol = "?protocol=3"
)
