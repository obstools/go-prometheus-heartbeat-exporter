package heartbeat

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewExporter(t *testing.T) {
	t.Run("creates new exporter", func(t *testing.T) {
		port, shutdownTimeout, route, logger := 8080, 42, "/some_route", newLogger(false, false)
		exporter := newExporter(port, shutdownTimeout, route, logger)

		assert.Equal(t, logger, exporter.logger)
		assert.Equal(t, time.Duration(shutdownTimeout), exporter.shutdownTimeout)
		assert.Equal(t, ":"+strconv.Itoa(port), exporter.server.Port())
		assert.Equal(t, route, exporter.route)
	})
}

func TestExporterStart(t *testing.T) {
	t.Run("starts exporter server", func(t *testing.T) {
		prometheusServer, logger := new(serverPrometheusMock), new(loggerMock)
		port, route, parentContext, wg := ":8080", "/some_route", context.Background(), new(sync.WaitGroup)
		exporter := &exporter{
			server: prometheusServer,
			logger: logger,
			route:  route,
		}
		prometheusServer.On("Port").Once().Return(port)
		logger.On("info", []string{exporterStartMessage + port + route}).Once().Return(nil)
		prometheusServer.On("ListenAndServe").Once().Return(nil)

		assert.NoError(t, exporter.start(parentContext, wg))
		assert.Equal(t, parentContext, exporter.ctx)
		assert.Equal(t, wg, exporter.wg)
		prometheusServer.AssertExpectations(t)
		logger.AssertExpectations(t)
	})
}

func TestExporterStop(t *testing.T) {
	t.Run("stops exporter server", func(t *testing.T) {
		prometheusServer := new(serverPrometheusMock)
		exporter := &exporter{
			server:          prometheusServer,
			ctx:             context.Background(),
			shutdownTimeout: time.Duration(2),
		}
		prometheusServer.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).Once().Return(nil)

		assert.NoError(t, exporter.stop())
		prometheusServer.AssertExpectations(t)
	})
}

func TestExporterListenShutdownSignal(t *testing.T) {
	prometheusServer, logger, wg := new(serverPrometheusMock), new(loggerMock), &sync.WaitGroup{}
	parentContext, cancel := context.WithCancel(context.Background())
	exporter := &exporter{
		server:          prometheusServer,
		logger:          logger,
		ctx:             parentContext,
		wg:              wg,
		shutdownTimeout: time.Duration(2),
	}
	t.Run("listens shutdown signal, exporter stops without error", func(t *testing.T) {
		wg.Add(1) // Add goroutine to the wait group
		logger.On("warning", []string{exporterShutdownMessage}).Once().Return(nil)
		logger.On("info", []string{exporterStopMessage}).Once().Return(nil)
		prometheusServer.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).Once().Return(nil)
		exporter.listenShutdownSignal()
		cancel()  // Simulate shutdown signal
		wg.Wait() // Wait for goroutine to finish

		assert.Nil(t, exporter.err)
		prometheusServer.AssertExpectations(t)
		logger.AssertExpectations(t)
	})

	t.Run("listens shutdown signal, exporter stops with error", func(t *testing.T) {
		wg.Add(1) // Add goroutine to the wait group
		logger.On("warning", []string{exporterShutdownMessage}).Once().Return(nil)
		logger.On("info", []string{exporterStopMessage}).Once().Return(nil)
		prometheusServer.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).Once().Return(errors.New("some error"))
		exporter.listenShutdownSignal()
		cancel()  // Simulate shutdown signal
		wg.Wait() // Wait for goroutine to finish

		assert.Error(t, exporter.err)
		prometheusServer.AssertExpectations(t)
		logger.AssertExpectations(t)
	})
}

func TestExporterIsPortAvailable(t *testing.T) {
	prometheusServer, port := new(serverPrometheusMock), ":8282"
	t.Run("checks if port is available", func(t *testing.T) {
		exporter := &exporter{
			server: prometheusServer,
			logger: new(loggerMock),
		}
		prometheusServer.On("Port").Once().Return(port)

		assert.NoError(t, exporter.isPortAvailable())
		prometheusServer.AssertExpectations(t)
	})

	t.Run("checks if port is unavailable", func(t *testing.T) {
		exporter := &exporter{
			server: prometheusServer,
			logger: new(loggerMock),
		}
		prometheusServer.On("Port").Once().Return(port)
		listener, _ := net.Listen("tcp", port)
		defer listener.Close()

		assert.EqualError(t, exporter.isPortAvailable(), exporterErrorMessage+port)
		prometheusServer.AssertExpectations(t)
	})
}
