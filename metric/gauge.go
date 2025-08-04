package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type GaugeVec struct {
	gauges *prometheus.GaugeVec
}

func NewGaugeVec(namespace, metricsName, help string, labels []string) *GaugeVec {
	cc := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      metricsName + "_g",
		Help:      help + " (gauges)",
	}, labels)

	prometheus.MustRegister(cc)

	return &GaugeVec{
		gauges: cc,
	}
}

func (v *GaugeVec) Inc(labels ...string) {
	v.gauges.WithLabelValues(labels...).Inc()
}

func (v *GaugeVec) Add(value float64, labels ...string) {
	v.gauges.WithLabelValues(labels...).Add(value)
}

func (v *GaugeVec) Dec(labels ...string) {
	v.gauges.WithLabelValues(labels...).Dec()
}

func (v *GaugeVec) Sub(value float64, labels ...string) {
	v.gauges.WithLabelValues(labels...).Sub(value)
}

func (v *GaugeVec) Set(value float64, labels ...string) {
	v.gauges.WithLabelValues(labels...).Set(value)
}
