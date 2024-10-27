package heartbeat

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Run("returns new server", func(t *testing.T) {
		shutdownTimeout, metricsRoute := 42, "/metrics"
		server := newServer(
			&Configuration{
				InstancesAttributes: []*InstanceAttributes{
					{
						Connection: "some_connection",
						URL:        "some_url",
					},
				},
				Port:            8080,
				ShutdownTimeout: shutdownTimeout,
				MetricsRoute:    metricsRoute,
				LogToStdout:     true,
				LogActivity:     true,
			},
		)

		assert.NotNil(t, server)
		assert.Equal(t, 1, len(server.heartbeatInstances))
		assert.Equal(t, time.Duration(shutdownTimeout), server.exporter.shutdownTimeout)
		assert.Equal(t, metricsRoute, server.exporter.route)
	})
}

func TestServerStart(t *testing.T) {
	t.Run("when no errors happens during starting and running the server", func(t *testing.T) {
		server := newServer(
			&Configuration{
				InstancesAttributes: []*InstanceAttributes{
					{
						Connection: "postgres",
						URL:        "postgres://localhost:5432/postgres",
					},
				},
				Port:         8080,
				MetricsRoute: "/metrics",
			},
		)

		assert.NoError(t, server.Start())
		assert.NotNil(t, server.ctx)
		assert.NotNil(t, server.shutdown)
		assert.True(t, server.isStarted())

		_ = server.Stop()
	})

	t.Run("when error happens during starting the server, server is already started", func(t *testing.T) {
		server, logger := &Server{started: true}, new(loggerMock)
		server.logger = logger
		logger.On("error", []string{serverStartErrorMessage}).Once()

		serverStart := server.Start()
		assert.Error(t, serverStart)
		assert.EqualError(t, serverStart, serverStartErrorMessage)
		logger.AssertExpectations(t)
	})

	t.Run("when error happens during starting the server, port is already in use", func(t *testing.T) {
		port := ":8080"
		listener, _ := net.Listen("tcp", port)
		defer listener.Close()
		time.Sleep(100 * time.Millisecond)
		server, logger, errMessage := newServer(createNewMinimalConfiguration()), new(loggerMock), exporterErrorMessage+port
		server.logger = logger
		logger.On("error", []string{errMessage}).Once()

		serverStart := server.Start()
		assert.Error(t, serverStart)
		assert.EqualError(t, serverStart, errMessage)
		assert.False(t, server.isStarted())
		logger.AssertExpectations(t)
	})
}

func TestServerStop(t *testing.T) {
	t.Run("when server is started, no errors happen during stopping exporter", func(t *testing.T) {
		server, wg, logger := newServer(createNewMinimalConfiguration()), new(waitGroupMock), new(loggerMock)
		server.wg, server.logger, server.started, server.shutdown = wg, logger, true, func() {}
		wg.On("Wait").Once()
		logger.On("info", []string{serverStopMessage}).Once()

		assert.NoError(t, server.Stop())
		assert.False(t, server.isStarted())
		wg.AssertExpectations(t)
		logger.AssertExpectations(t)
	})

	t.Run("when server is started, exporter returns error during stopping", func(t *testing.T) {
		server, wg, logger, err := newServer(createNewMinimalConfiguration()), new(waitGroupMock), new(loggerMock), errors.New("some error")
		server.wg, server.logger, server.started, server.shutdown, server.exporter.err = wg, logger, true, func() {}, err
		wg.On("Wait").Once()
		logger.On("info", []string{serverStopMessage}).Once()
		logger.On("warning", []string{serverStopExporterErrorMessage, err.Error()}).Once()

		serverStop := server.Stop()
		assert.Error(t, serverStop)
		assert.EqualError(t, serverStop, err.Error())
		wg.AssertExpectations(t)
		logger.AssertExpectations(t)
	})

	t.Run("when server is not started", func(t *testing.T) {
		server, logger := new(Server), new(loggerMock)
		server.logger = logger
		logger.On("error", []string{serverStopErrorMessage}).Once()

		serverStop := server.Stop()
		assert.Error(t, serverStop)
		assert.EqualError(t, serverStop, serverStopErrorMessage)
		logger.AssertExpectations(t)
	})
}

func TestServerIsStarted(t *testing.T) {
	t.Run("when server is started", func(t *testing.T) {
		server := new(Server)
		server.started = true

		assert.True(t, server.isStarted())
	})

	t.Run("when server is not started", func(t *testing.T) {
		server := new(Server)

		assert.False(t, server.isStarted())
	})
}

func TestServerStartPrivate(t *testing.T) {
	server := new(Server)
	server.start()

	assert.True(t, server.started)
}

func TestServerStopPrivate(t *testing.T) {
	server := &Server{started: true}
	server.stop()

	assert.False(t, server.started)
}
