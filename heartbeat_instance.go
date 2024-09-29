package heartbeat

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Heartbeat instance structure
type heartbeatInstance struct {
	name                    string
	intervalSec, timeoutSec int
	metric                  metric
	session                 session
	ctx                     context.Context
	wg                      *sync.WaitGroup
	logger                  logger
}

// New heartbeatInstance builder. Returns new instance pointer
func newInstance(instanceAttrs *InstanceAttributes) *heartbeatInstance {
	return &heartbeatInstance{
		name:        instanceAttrs.Name,
		intervalSec: instanceAttrs.IntervalSec,
		timeoutSec:  instanceAttrs.TimeoutSec,
		session:     newSession(instanceAttrs.Connection, instanceAttrs.URL),
	}
}

// heartbeatInstance methods

// Heartbeat instance wrapper with execution timeout limitter
func (instance *heartbeatInstance) heartbeat(ctx context.Context, successChannel chan float64, failureChannel chan *heartbeatError) {
	defer close(successChannel)
	defer close(failureChannel)

	// Runs heartbeat as instance session
	startTime, err := time.Now(), instance.session.run()
	if err != nil {
		failureChannel <- &heartbeatError{duration: time.Since(startTime).Seconds(), err: err}
		ctx.Done()
		return
	}

	successChannel <- time.Since(startTime).Seconds()
}

// Heartbeat instance wrapper with result metrics
func (instance *heartbeatInstance) heartbeatRunner() {
	ctx, cancel := context.WithTimeout(instance.ctx, time.Duration(instance.timeoutSec)*time.Second)
	defer cancel()

	successChannel, failureChannel := make(chan float64, 1), make(chan *heartbeatError)
	name, metric := instance.name, instance.metric
	go instance.heartbeat(ctx, successChannel, failureChannel)

	select {
	case duration := <-successChannel:
		metric.setSuccess(duration)
		instance.logger.infoActivity(
			heartbeatActivitySuccess,
			name,
			heartbeatActivitySuccessElapsedTime,
			fmt.Sprintf("%.2f", duration),
		)
	case failedHealthcheck := <-failureChannel:
		metric.setFailureConnection()
		instance.logger.infoActivity(
			heartbeatActivityFailure,
			name,
			heartbeatActivityFailureSessionError,
			failedHealthcheck.originalError().Error(),
		)
	case <-ctx.Done():
		metric.setFailureTimeout()
		instance.logger.infoActivity(
			heartbeatActivityFailure,
			name,
			heartbeatActivityFailureSessionError,
			ctx.Err().Error(), // TODO: Good point to do this error message to be configurable
		)
	}
}

// Instance heartbeat metric builder. Sets metric pointer to instance field
func (instance *heartbeatInstance) buildHeartbeatMetric() {
	instance.metric = newMetric(instance.name)
}

// Heartbeat instance worker. Runs loop with target instance heartbeat
func (instance *heartbeatInstance) workerRunner() {
	defer instance.wg.Done()

	workerName := instance.name
	instance.buildHeartbeatMetric()
	instance.logger.info(heartbeatWorkerMsg, workerName, heartbeatWorkerStatusActive)

	for {
		select {
		case <-instance.ctx.Done():
			instance.logger.info(heartbeatWorkerMsg, workerName, heartbeatWorkerStatusInactive)
			return
		default:
			instance.heartbeatRunner()
			time.Sleep(time.Duration(instance.intervalSec) * time.Second)
		}
	}
}
