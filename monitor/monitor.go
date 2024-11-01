package monitor

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Interface interface {
	PathStatusCode(path string, code int)
	StatusCacheHit()
	StatusCacheMiss()
	StatusDurationStart(path string) *prometheus.Timer
	StatusDurationEnd(timer *prometheus.Timer)
}

type Monitor struct {
	pathStatusCodes     *prometheus.CounterVec
	pathStatusCacheHit  prometheus.Counter
	pathStatusCacheMiss prometheus.Counter
	pathStatusDuration  *prometheus.HistogramVec
}

func New() *Monitor {
	pathStatusCodes := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_path_status_codes",
		Help: "The http statuc code per path",
	}, []string{"path", "status_code"})

	pathStatusCacheHit := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "status_cache_hit",
		Help: "The number of status cache hits",
	})

	pathStatusCacheMiss := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "status_cache_miss",
		Help: "The number of status cache misses",
	})

	pathStatusDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "status_http_duration",
		Help: "The status of the status endpoint in seconds",
	}, []string{"path"})

	prometheus.Register(pathStatusCodes)
	prometheus.Register(pathStatusCacheHit)
	prometheus.Register(pathStatusCacheMiss)
	prometheus.Register(pathStatusDuration)

	return &Monitor{
		pathStatusCodes:     pathStatusCodes,
		pathStatusCacheHit:  pathStatusCacheHit,
		pathStatusCacheMiss: pathStatusCacheMiss,
		pathStatusDuration:  pathStatusDuration,
	}
}

func (m *Monitor) PathStatusCode(path string, code int) {
	m.pathStatusCodes.WithLabelValues(path, strconv.Itoa(code)).Inc()
}

func (m *Monitor) StatusCacheHit() {
	m.pathStatusCacheHit.Inc()
}

func (m *Monitor) StatusCacheMiss() {
	m.pathStatusCacheMiss.Inc()
}

func (m *Monitor) StatusDurationStart(path string) *prometheus.Timer {
	return prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		m.pathStatusDuration.WithLabelValues(path).Observe(v)
	}))
}

func (m *Monitor) StatusDurationEnd(timer *prometheus.Timer) {
	timer.ObserveDuration()
}
