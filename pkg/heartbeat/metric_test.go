package heartbeat

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	t.Run("returns new metric", func(t *testing.T) {
		metric := newMetric(generateUniqueInstanceName())

		assert.NotNil(t, metric)
	})
}

func TestInstanceMetricSetSuccess(t *testing.T) {
	t.Run("sets success metric", func(t *testing.T) {
		prometheusCounterMock, prometheusGaugeMock, duration := &prometheusCounterMock{}, &prometheusGaugeMock{}, 42.0
		metric := &instanceMetric{success: prometheusCounterMock, duration: prometheusGaugeMock}
		prometheusCounterMock.On("Inc").Return()
		prometheusGaugeMock.On("Set", duration).Return()
		metric.setSuccess(duration)

		prometheusCounterMock.AssertExpectations(t)
		prometheusGaugeMock.AssertExpectations(t)
	})
}

func TestInstanceMetricSetFailureMetric(t *testing.T) {
	t.Run("sets failure metric", func(t *testing.T) {
		label, prometheusCounterMock, prometheusCounterVecMock := "some label", &prometheusCounterMock{}, &prometheusCounterVecMock{}
		metric := &instanceMetric{failure: prometheusCounterVecMock}
		prometheusCounterVecMock.On("WithLabelValues", []string{label}).Return(prometheusCounterMock)
		prometheusCounterMock.On("Add", 1.0).Return()
		metric.setFailureMetric(label)

		prometheusCounterVecMock.AssertExpectations(t)
		prometheusCounterMock.AssertExpectations(t)
	})
}

func TestInstanceMetricSetFailureConnection(t *testing.T) {
	t.Run("sets failure connection metric", func(t *testing.T) {
		prometheusCounterMock, prometheusCounterVecMock := &prometheusCounterMock{}, &prometheusCounterVecMock{}
		metric := &instanceMetric{failure: prometheusCounterVecMock}
		prometheusCounterVecMock.On("WithLabelValues", []string{metricFailureLableConncection}).Return(prometheusCounterMock)
		prometheusCounterMock.On("Add", 1.0).Return()
		metric.setFailureConnection()

		prometheusCounterVecMock.AssertExpectations(t)
		prometheusCounterMock.AssertExpectations(t)
	})
}

func TestInstanceMetricSetFailureTimeout(t *testing.T) {
	t.Run("sets failure timeout metric", func(t *testing.T) {
		prometheusCounterMock, prometheusCounterVecMock := &prometheusCounterMock{}, &prometheusCounterVecMock{}
		metric := &instanceMetric{failure: prometheusCounterVecMock}
		prometheusCounterVecMock.On("WithLabelValues", []string{metricFailureLableTimeout}).Return(prometheusCounterMock)
		prometheusCounterMock.On("Add", 1.0).Return()
		metric.setFailureTimeout()

		prometheusCounterVecMock.AssertExpectations(t)
		prometheusCounterMock.AssertExpectations(t)
	})
}

func TestPrometheusCounterVecWrapperWithLabelValues(t *testing.T) {
	t.Run("returns prometheus counter with label values", func(t *testing.T) {
		counterVec := promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        metricNameFailedAttempts,
				Help:        metricDescFailedAttempts,
				ConstLabels: prometheus.Labels{metricLabelInstanceName: generateUniqueInstanceName()},
			},
			[]string{metricLabelErrorType},
		)
		wrapper := &prometheusCounterVecWrapper{CounterVec: counterVec}

		assert.NotNil(t, wrapper.WithLabelValues("some label"))
	})
}
