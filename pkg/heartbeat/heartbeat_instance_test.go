package heartbeat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewInstance(t *testing.T) {
	name, url, intervalSec, timeoutSec := generateUniqueInstanceName(), "test_url", 1, 2
	t.Run("when connection is postgres", func(t *testing.T) {
		connection := "postgres"
		instanceAttributes := &InstanceAttributes{
			Connection:  connection,
			Name:        name,
			URL:         url,
			IntervalSec: intervalSec,
			TimeoutSec:  timeoutSec,
		}
		instance := newInstance(instanceAttributes)
		instanceSession := instance.session

		assert.Equal(t, instanceAttributes.Name, instance.name)
		assert.Equal(t, instanceAttributes.IntervalSec, instance.intervalSec)
		assert.Equal(t, instanceAttributes.TimeoutSec, instance.timeoutSec)
		assert.NotNil(t, instanceSession)
		assert.Equal(t, connection, instanceSession.getConnection())
		assert.Equal(t, url, instanceSession.getURL())
	})

	t.Run("when connection is undefined", func(t *testing.T) {
		connection := "undefined"
		instanceAttributes := &InstanceAttributes{
			Connection:  connection,
			Name:        name,
			URL:         url,
			IntervalSec: intervalSec,
			TimeoutSec:  timeoutSec,
		}
		instance := newInstance(instanceAttributes)

		assert.Equal(t, instanceAttributes.Name, instance.name)
		assert.Equal(t, instanceAttributes.IntervalSec, instance.intervalSec)
		assert.Equal(t, instanceAttributes.TimeoutSec, instance.timeoutSec)
		assert.Nil(t, instance.session)
	})
}

func TestHeartbeatInstanceHeartbeat(t *testing.T) {
	t.Run("when session does not return error", func(t *testing.T) {
		ctx, successChannel, failureChannel := context.Background(), make(chan float64, 1), make(chan *heartbeatError, 1)
		session := new(sessionMock)
		instance := newTestInstance(session)
		session.On("run").Once().Return(nil)
		go instance.heartbeat(ctx, successChannel, failureChannel)

		assert.Greater(t, <-successChannel, 0.0)
		assert.Zero(t, len(successChannel))
		assert.Zero(t, len(failureChannel))
		session.AssertExpectations(t)
	})

	t.Run("when session returns error", func(t *testing.T) {
		ctx, successChannel, failureChannel := context.Background(), make(chan float64, 1), make(chan *heartbeatError, 1)
		session, err := new(sessionMock), errors.New("test error")
		instance := newTestInstance(session)
		session.On("run").Once().Return(err)
		go instance.heartbeat(ctx, successChannel, failureChannel)
		heartbeatError := <-failureChannel

		assert.Zero(t, len(failureChannel))
		assert.NotZero(t, heartbeatError.duration)
		assert.Equal(t, heartbeatError.err, err)
		assert.Zero(t, len(successChannel))
		session.AssertExpectations(t)
	})
}

func TestHeartbeatInstanceHeartbeatRunner(t *testing.T) {
	t.Run("when instance heartbeat success", func(t *testing.T) {
		var capturedDuration float64
		session, metric, logger := new(sessionMock), new(metricMock), new(loggerMock)
		instance := newTestInstance(session)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		instance.ctx, instance.wg, instance.metric, instance.logger = ctx, createNewWaitGroup(), metric, logger
		instance.wg.Add(1)
		session.On("run").Once().Return(nil)
		metric.On("setSuccess", mock.AnythingOfType("float64")).Run(func(args mock.Arguments) {
			capturedDuration = args.Get(0).(float64)
		}).Once()
		logger.On("infoActivity", mock.MatchedBy(func(args []string) bool {
			return len(args) == 4 &&
				args[0] == heartbeatActivitySuccess &&
				args[1] == instance.name &&
				args[2] == heartbeatActivitySuccessElapsedTime &&
				len(args[3]) > 0
		})).Once().Return(nil)
		instance.heartbeatRunner()

		assert.Greater(t, capturedDuration, 0.0, "Duration should be greater than 0")
		assert.Less(t, capturedDuration, float64(instance.timeoutSec), "Duration should be less than timeout")
		session.AssertExpectations(t)
		metric.AssertExpectations(t)
		logger.AssertExpectations(t)
	})

	t.Run("when instance heartbeat failure, connection error", func(t *testing.T) {
		session, metric, logger, errorMessage := new(sessionMock), new(metricMock), new(loggerMock), "test error"
		instance, err := newTestInstance(session), errors.New(errorMessage)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		instance.ctx, instance.wg, instance.metric, instance.logger = ctx, createNewWaitGroup(), metric, logger
		instance.wg.Add(1)
		session.On("run").Once().Return(err)
		metric.On("setFailureConnection").Once()
		logger.On("infoActivity", []string{
			heartbeatActivityFailure,
			instance.name,
			heartbeatActivityFailureSessionError,
			errorMessage,
		}).Once()
		instance.heartbeatRunner()

		session.AssertExpectations(t)
		metric.AssertExpectations(t)
		logger.AssertExpectations(t)
	})

	t.Run("when instance heartbeat failure, timeout error", func(t *testing.T) {
		session, metric, logger := new(sessionMock), new(metricMock), new(loggerMock)
		instance := newTestInstance(session)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Millisecond)
		defer cancel()

		instance.ctx, instance.wg, instance.metric, instance.logger = ctx, createNewWaitGroup(), metric, logger
		instance.wg.Add(1)

		// Simulate a long-running operation that will cause a timeout
		session.On("run").Once().Run(func(_ mock.Arguments) {
			time.Sleep(time.Duration(5) * time.Millisecond)
		}).Return(nil)

		metric.On("setFailureTimeout").Once()
		logger.On("infoActivity", []string{
			heartbeatActivityFailure,
			instance.name,
			heartbeatActivityFailureSessionError,
			"context deadline exceeded",
		}).Once()
		instance.heartbeatRunner()

		// Wait a bit to ensure all goroutines have finished
		time.Sleep(time.Duration(10) * time.Millisecond)

		session.AssertExpectations(t)
		metric.AssertExpectations(t)
		logger.AssertExpectations(t)
	})
}

func TestHeartbeatInstanceBuildHeartbeatMetric(t *testing.T) {
	instance := new(heartbeatInstance)
	instance.name = generateUniqueInstanceName()
	instance.buildHeartbeatMetric()

	assert.NotNil(t, instance.metric)
}

func TestHeartbeatInstanceWorkerRunner(t *testing.T) {
	t.Run("when instance heartbeat success and then context is cancelled", func(t *testing.T) {
		session, metric, logger := new(sessionMock), new(metricMock), new(loggerMock)
		instance := newTestInstance(session)

		ctx, cancel := context.WithCancel(context.Background())

		instance.ctx, instance.wg, instance.metric, instance.logger = ctx, createNewWaitGroup(), metric, logger
		instance.wg.Add(1)
		session.On("run").Once().Return(nil)
		logger.On("info", []string{
			heartbeatWorkerMsg,
			instance.name,
			heartbeatWorkerStatusActive,
		}).Once()
		logger.On("infoActivity", mock.MatchedBy(func(args []string) bool {
			return len(args) == 4 &&
				args[0] == heartbeatActivitySuccess &&
				args[1] == instance.name &&
				args[2] == heartbeatActivitySuccessElapsedTime &&
				len(args[3]) > 0
		})).Once()
		logger.On("info", []string{
			heartbeatWorkerMsg,
			instance.name,
			heartbeatWorkerStatusInactive,
		}).Once()

		go func() {
			time.Sleep(10 * time.Millisecond) // Give some time for heartbeatRunner to execute
			cancel()
		}()
		instance.workerRunner()
		instance.wg.Wait() // Wait for the goroutine to finish

		session.AssertExpectations(t)
		logger.AssertExpectations(t)
		metric.AssertExpectations(t)
	})
}
