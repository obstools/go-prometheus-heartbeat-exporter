package heartbeat

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metric interface
type metric interface {
	setSuccess(duration float64)
	setFailureConnection()
	setFailureTimeout()
}

// prometheusCounter interface
type prometheusCounter interface {
	Inc()
	Add(float64)
}

// prometheusCounterVec interface
type prometheusCounterVec interface {
	WithLabelValues(...string) prometheusCounter
}

// prometheusGauge interface
type prometheusGauge interface {
	Set(val float64)
}

// Instance Prometheus metric structure
type instanceMetric struct {
	success  prometheusCounter
	failure  prometheusCounterVec
	duration prometheusGauge
}

// prometheusCounterVecWrapper type to implement prometheusCounter interface
type prometheusCounterVecWrapper struct {
	*prometheus.CounterVec
}

// Prometheus metric structure builder. Returns new metric pointer
func newMetric(instanceName string) *instanceMetric {
	instanceLable, additionalLabels := prometheus.Labels{metricLabelInstanceName: instanceName}, []string{metricLabelErrorType}

	return &instanceMetric{
		success: promauto.NewCounter(prometheus.CounterOpts{
			Name:        metricNameSuccessfulAttempts,
			Help:        metricDescSuccessfulAttempts,
			ConstLabels: instanceLable,
		}),
		failure: &prometheusCounterVecWrapper{promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        metricNameFailedAttempts,
				Help:        metricDescFailedAttempts,
				ConstLabels: instanceLable,
			},
			additionalLabels,
		)},
		duration: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        metricNameDuration,
			Help:        metricDescDuration,
			ConstLabels: instanceLable,
		}),
	}
}

// prometheusCounterVecWrapper methods
// WithLabelValues method to implement prometheusCounter interface
func (wrapper *prometheusCounterVecWrapper) WithLabelValues(args ...string) prometheusCounter {
	return wrapper.CounterVec.WithLabelValues(args...)
}

// instanceMetric methods
// Sets successful metric with duration
func (instanceMetric *instanceMetric) setSuccess(duration float64) {
	instanceMetric.success.Inc()
	instanceMetric.duration.Set(duration)
}

// Sets failed metric
func (instanceMetric *instanceMetric) setFailureMetric(label string) {
	instanceMetric.failure.WithLabelValues(label).Add(1)
}

// Sets failed metric with connection label
func (instanceMetric *instanceMetric) setFailureConnection() {
	instanceMetric.setFailureMetric(metricFailureLableConncection)
}

// Sets failed metric with timeout label
func (instanceMetric *instanceMetric) setFailureTimeout() {
	instanceMetric.setFailureMetric(metricFailureLableTimeout)
}
