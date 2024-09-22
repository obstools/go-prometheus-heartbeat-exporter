package heartbeat

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metric structure
type metric struct {
	success  prometheus.Counter
	failure  *prometheus.CounterVec
	duration prometheus.Gauge
}

// Prometheus metric structure builder. Returns new metric pointer
func newMetric(instanceName string) *metric {
	instanceLable, additionalLabels := prometheus.Labels{metricLabelInstanceName: instanceName}, []string{metricLabelErrorType}

	return &metric{
		success: promauto.NewCounter(prometheus.CounterOpts{
			Name:        metricNameSuccessfulAttempts,
			Help:        metricDescSuccessfulAttempts,
			ConstLabels: instanceLable,
		}),
		failure: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name:        metricNameFailedAttempts,
				Help:        metricDescFailedAttempts,
				ConstLabels: instanceLable,
			},
			additionalLabels,
		),
		duration: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        metricNameDuration,
			Help:        metricDescDuration,
			ConstLabels: instanceLable,
		}),
	}
}
