package heartbeat

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// WaitGroup interface
type waitGroup interface {
	Add(int)
	Done()
	Wait()
}

// serverPrometheusWrapper structure. Used for testing purposes
type serverPrometheusWrapper struct {
	*http.Server
}

// serverPrometheusWrapper methods
func (wrapper *serverPrometheusWrapper) Port() string {
	return wrapper.Addr
}

// serverPrometheus interface
type serverPrometheus interface {
	ListenAndServe() error
	Shutdown(context.Context) error
	Port() string
}

// Exporter structure
type exporter struct {
	server          serverPrometheus
	shutdownTimeout time.Duration
	err             error
	route           string
	ctx             context.Context
	wg              waitGroup
	logger          logger
}

// Exporter builder. Returns pointer to new exporter structure
func newExporter(port, shutdownTimeout int, route string, logger logger) *exporter {
	handler := http.NewServeMux()
	handler.Handle(route, promhttp.Handler())

	return &exporter{
		server: &serverPrometheusWrapper{
			Server: &http.Server{
				Addr:    ":" + strconv.Itoa(port),
				Handler: handler,
			},
		},
		logger:          logger,
		shutdownTimeout: time.Duration(shutdownTimeout),
		route:           route,
	}
}

// Exporter methods

// Starts exporter, runs listen channel from the parent (heartbeat server)
func (exporter *exporter) start(parentContext context.Context, wg *sync.WaitGroup) error {
	exporter.ctx, exporter.wg = parentContext, wg
	exporter.listenShutdownSignal()
	exporter.logger.info(exporterStartMessage + exporter.server.Port() + exporter.route)

	return exporter.server.ListenAndServe()
}

// Stops exporter with timeout
func (exporter *exporter) stop() error {
	ctx, cancel := context.WithTimeout(exporter.ctx, exporter.shutdownTimeout*time.Second)
	defer cancel()

	return exporter.server.Shutdown(ctx)
}

// Exporter shutdown signal listener. Stops exporter by shutdown signal
func (exporter *exporter) listenShutdownSignal() {
	go func() {
		defer exporter.wg.Done()
		exporterLogger := exporter.logger

		<-exporter.ctx.Done()
		exporterLogger.warning(exporterShutdownMessage)
		if err := exporter.stop(); err != nil {
			exporter.err = err
		}

		exporterLogger.info(exporterStopMessage)
	}()
}

// Exporter port-for-bind checker
func (exporter *exporter) isPortAvailable() (err error) {
	port := exporter.server.Port()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return errors.New(exporterErrorMsg + port)
	}

	listener.Close()
	return
}
