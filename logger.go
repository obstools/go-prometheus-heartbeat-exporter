package heartbeat

import (
	"io"
	"log"
	"os"
)

// TODO: uncomment when this logger will embedded as interface
// // Logger interface
// type logger interface {
// 	infoActivity(...string)
// 	info(...string)
// 	warning(...string)
// 	error(...string)
// }

// Custom logger that supports 3 different log levels (info, warning, error)
type eventLogger struct {
	eventInfo, eventWarning, eventError *log.Logger
	logToStdout, logActivity            bool
	flag                                int
	stdout, stderr                      io.Writer
}

// Logger builder. Returns pointer to new logger structure
func newLogger(logToStdout, logActivity bool) *eventLogger {
	return &eventLogger{
		logToStdout: logToStdout,
		logActivity: logActivity,
		flag:        logFlag,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
	}
}

// logger methods

// Provides INFO log level for server activities. Writes to stdout for case when
// logger.logToStdout and logger.logActivity are enabled, suppressed otherwise
func (logger *eventLogger) infoActivity(message ...string) {
	if logger.logToStdout && logger.logActivity {
		if logger.eventInfo == nil {
			logger.eventInfo = log.New(logger.stdout, infoLogLevel+": ", logger.flag)
		}

		logger.eventInfo.Println(message)
	}
}

// Provides INFO log level. Writes to stdout for case when logger.logToStdout is enabled,
// suppressed otherwise
func (logger *eventLogger) info(message ...string) {
	if logger.logToStdout {
		if logger.eventInfo == nil {
			logger.eventInfo = log.New(logger.stdout, infoLogLevel+": ", logger.flag)
		}

		logger.eventInfo.Println(message)
	}
}

// Provides WARNING log level. Writes to stdout for case when logger.logToStdout is enabled,
// suppressed otherwise
func (logger *eventLogger) warning(message ...string) {
	if logger.logToStdout {
		if logger.eventWarning == nil {
			logger.eventWarning = log.New(logger.stdout, warningLogLevel+": ", logger.flag)
		}

		logger.eventWarning.Println(message)
	}
}

// Provides ERROR log level. Writes to stdout for case when logger.logToStdout is enabled,
// suppressed otherwise
func (logger *eventLogger) error(message ...string) {
	if logger.logToStdout {
		if logger.eventError == nil {
			logger.eventError = log.New(logger.stderr, errorLogLevel+": ", logger.flag)
		}

		logger.eventError.Println(message)
	}
}
