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
