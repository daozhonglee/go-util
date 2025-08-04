package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CounterVec struct {
	counters *prometheus.CounterVec
}

func NewCounterVec(namespace, metricsName, help string, labels []string) *CounterVec {
	cc := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      metricsName + "_c",
		Help:      help + " (counters)",
	}, labels)

	prometheus.MustRegister(cc)

	return &CounterVec{
		counters: cc,
	}
}

func (v *CounterVec) Inc(labels ...string) {
	v.counters.WithLabelValues(labels...).Inc()
}

func (v *CounterVec) Add(count float64, labels ...string) {
	v.counters.WithLabelValues(labels...).Add(count)
}
