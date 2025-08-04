package metric

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/daozhonglee/go-util/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cv *CounterVec
var gv *GaugeVec
var grpcTimer *Timer
var httpTimer *Timer
var customTimer *Timer

type Server struct {
	lsnAddr string
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%s", s.lsnAddr))
	if err != nil {
		log.CRITICAL("[MetricServer] start err: %v, lsn: %v", err, s.lsnAddr)
		return
	}

	http.Handle("/debug/metrics", promhttp.Handler())
	err = http.Serve(listener, nil)
	if err != nil {
		log.CRITICAL("[MetricServer] Serve err: %v, lsn: %v", err, s.lsnAddr)
		return
	}
	log.INFO("metric server start success!")
}

func initBaseMetric(module string) {
	cv = NewCounterVec(NameSpaceSugo, fmt.Sprintf("%s_%s", module, "counter"), "counter_vec", []string{"target", "detail", "result"})
	gv = NewGaugeVec(NameSpaceSugo, fmt.Sprintf("%s_%s", module, "gauge"), "gauge_vec", []string{"target", "detail", "result"})
	customTimer = NewTimer(NameSpaceSugo, fmt.Sprintf("%s_%s", module, "timer"), "custom_timer", []string{"action", "result"})
}

// Init InitGRPC，初始化GRPC
func Init(listen string, module string) {
	server := &Server{lsnAddr: listen}
	go server.Start()

	initBaseMetric(module)
	grpcTimer = NewTimer(NameSpaceSugo, fmt.Sprintf("%s_%s", module, "grpc"), "grpc_timer", []string{"action", "result"})
	log.DEBUG("[Metric] Init listen: %s, module: %s", listen, module)
}

// InitHttp 初始化Http
func InitHttp(listen string, module string) {
	server := &Server{lsnAddr: listen}
	go server.Start()

	initBaseMetric(module)
	httpTimer = NewTimer(NameSpaceSugo, fmt.Sprintf("%s_%s", module, "http"), "http_timer", []string{"method", "uri", "code"})
	log.DEBUG("[Metric] InitHttp listen: %s, module: %s", listen, module)
}

func InitJob(listen string, module string) {
	server := &Server{lsnAddr: listen}
	go server.Start()

	initBaseMetric(module)
	log.DEBUG("[Metric] InitJob listen: %s, module: %s", listen, module)
}

func IncCounter(labels ...string) {
	if cv != nil {
		cv.counters.WithLabelValues(labels...).Inc()
	}
}

func AddCounter(count float64, labels ...string) {
	if cv != nil {
		cv.counters.WithLabelValues(labels...).Add(count)
	}
}

func BeginTimerGRPC() func(values ...string) {
	if grpcTimer != nil {
		return grpcTimer.Timer()
	}
	return nil
}

func BeginTimerHttp() func(values ...string) {
	if httpTimer != nil {
		return httpTimer.Timer()
	}
	return nil
}

func BeginTimerCustom() func(values ...string) {
	if customTimer != nil {
		return customTimer.Timer()
	}
	return nil
}

func EndTimer(tf func(values ...string), action string, err error) {
	if tf != nil {
		if err != nil {
			tf(action, "fail")
		} else {
			tf(action, "success")
		}
	}
}

func EndTimerHttp(tf func(values ...string), method string, uri string, code int) {
	if tf != nil {
		tf(method, uri, strconv.Itoa(code))
	}
}

func EndTimerWithCode(tf func(values ...string), action string, code int) {
	if tf != nil {
		tf(action, strconv.Itoa(code))
	}
}

func IncGauge(labels ...string) {
	if gv != nil {
		gv.gauges.WithLabelValues(labels...).Inc()
	}
}

func AddGauge(value float64, labels ...string) {
	if gv != nil {
		gv.gauges.WithLabelValues(labels...).Add(value)
	}
}

func DecGauge(labels ...string) {
	if gv != nil {
		gv.gauges.WithLabelValues(labels...).Dec()
	}
}

func SubGauge(value float64, labels ...string) {
	if gv != nil {
		gv.gauges.WithLabelValues(labels...).Sub(value)
	}
}

func SetGauge(value float64, labels ...string) {
	if gv != nil {
		gv.gauges.WithLabelValues(labels...).Set(value)
	}
}
