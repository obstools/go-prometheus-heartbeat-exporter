package heartbeat

import (
	"context"
	"errors"
	"sync"
)

// WaitGroup interface
type waitGroup interface {
	Add(int)
	Done()
	Wait()
}

// Server structure
type Server struct {
	heartbeatInstances []*heartbeatInstance
	logger             logger
	ctx                context.Context
	shutdown           context.CancelFunc
	wg                 waitGroup
	started            bool
	sync.Mutex

	exporter *exporter
}

// Server builder. Returns pointer to new server structure
func newServer(configuration *Configuration) *Server {
	var heartbeatInstances []*heartbeatInstance

	for _, instanceAttributes := range configuration.InstancesAttributes {
		heartbeatInstances = append(heartbeatInstances, newInstance(instanceAttributes))
	}
	logger := newLogger(configuration.LogToStdout, configuration.LogActivity)

	return &Server{
		heartbeatInstances: heartbeatInstances,
		logger:             logger,
		wg:                 new(sync.WaitGroup),
		exporter: newExporter(
			configuration.Port,
			configuration.ShutdownTimeout,
			configuration.MetricsRoute,
			logger,
		),
	}
}

// Server methods

// Starts server. Returns error if any
func (server *Server) Start() (err error) {
	logger := server.logger

	if server.isStarted() {
		err = errors.New(serverStartErrorMessage)
		logger.error(err.Error())

		return err
	} else if err := server.exporter.isPortAvailable(); err != nil {
		logger.error(err.Error())

		return err
	}

	logger.info(serverStartMessage)
	server.ctx, server.shutdown = context.WithCancel(context.Background())

	server.wg.Add(1)
	go func() {
		// We have checked port availability before, it's safe to start exporter
		if err := server.exporter.start(server.ctx, server.wg); err != nil {
			logger.info(exporterStatusTitle, err.Error())
		}
	}()

	for _, instance := range server.heartbeatInstances {
		instance.ctx, instance.wg, instance.logger = server.ctx, server.wg, server.logger
		server.wg.Add(1)
		go instance.workerRunner()
	}
	server.start()

	return err
}

// Stops server. Returns error if server is not started
func (server *Server) Stop() (err error) {
	logger := server.logger

	if server.isStarted() {
		server.shutdown()
		server.wg.Wait()
		server.stop()
		logger.info(serverStopMessage)
		if err = server.exporter.err; err != nil {
			logger.warning(serverStopExporterErrorMessage, err.Error())
		}

		return err
	}

	err = errors.New(serverStopErrorMessage)
	logger.error(err.Error())

	return err
}

// Thread-safe getter to check if server has been started. Returns server.started
func (server *Server) isStarted() bool {
	server.Lock()
	defer server.Unlock()
	return server.started
}

// Thread-safe setter of started-flag to indicate server has been started
func (server *Server) start() {
	server.Lock()
	defer server.Unlock()
	server.started = true
}

// Thread-safe setter of started-flag to indicate server has been stopped
func (server *Server) stop() {
	server.Lock()
	defer server.Unlock()
	server.started = false
}
