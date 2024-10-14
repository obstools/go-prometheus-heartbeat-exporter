package heartbeat

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Testing mocks

// logger mock
type loggerMock struct {
	mock.Mock
}

func (mock *loggerMock) infoActivity(message ...string) {
	mock.Called(message)
}

func (mock *loggerMock) info(message ...string) {
	mock.Called(message)
}

func (mock *loggerMock) warning(message ...string) {
	mock.Called(message)
}

func (mock *loggerMock) error(message ...string) {
	mock.Called(message)
}

// serverPrometheus mock
type serverPrometheusMock struct {
	mock.Mock
}

func (mock *serverPrometheusMock) ListenAndServe() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *serverPrometheusMock) Shutdown(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

func (mock *serverPrometheusMock) Port() string {
	args := mock.Called()
	return args.String(0)
}

// session mock
type sessionMock struct {
	mock.Mock
}

func (mock *sessionMock) run() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *sessionMock) getConnection() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *sessionMock) getURL() string {
	args := mock.Called()
	return args.String(0)
}

// prometheusCounterMock
type prometheusCounterMock struct {
	mock.Mock
}

func (mock *prometheusCounterMock) Inc() {
	mock.Called()
}

func (mock *prometheusCounterMock) Add(val float64) {
	mock.Called(val)
}

// prometheusCounterVecMock
type prometheusCounterVecMock struct {
	mock.Mock
}

func (mock *prometheusCounterVecMock) WithLabelValues(labels ...string) prometheusCounter {
	args := mock.Called(labels)
	return args.Get(0).(prometheusCounter)
}

// prometheusGaugeMock
type prometheusGaugeMock struct {
	mock.Mock
}

func (mock *prometheusGaugeMock) Set(val float64) {
	mock.Called(val)
}

// metricMock
type metricMock struct {
	mock.Mock
}

func (mock *metricMock) setSuccess(duration float64) {
	mock.Called(duration)
}

func (mock *metricMock) setFailureConnection() {
	mock.Called()
}

func (mock *metricMock) setFailureTimeout() {
	mock.Called()
}

// // WaitGroup mock
// type waitGroupMock struct {
// 	mock.Mock
// }

// func (mock *waitGroupMock) Add(count int) {
// 	mock.Called(count)
// }

// func (mock *waitGroupMock) Done() {
// 	mock.Called()
// }

// func (mock *waitGroupMock) Wait() {
// 	mock.Called()
// }
