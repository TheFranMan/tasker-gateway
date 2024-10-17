package monitor

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Interface interface {
	PathStatusCode(path string, code int)
	PathStatusCached()
}

type Monitor struct {
	pathStatusCodes  *prometheus.CounterVec
	pathStatusCached prometheus.Counter
}

func New() *Monitor {
	pathStatusCodes := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_path_status_codes",
		Help: "The http statuc code per path",
	}, []string{"path", "status_code"})

	pathStatusCached := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_status_cache_hit",
		Help: "The number of status cache hits",
	})

	prometheus.Register(pathStatusCodes)
	prometheus.Register(pathStatusCached)

	return &Monitor{
		pathStatusCodes:  pathStatusCodes,
		pathStatusCached: pathStatusCached,
	}
}

func (m *Monitor) PathStatusCode(path string, code int) {
	m.pathStatusCodes.WithLabelValues(path, strconv.Itoa(code)).Inc()
}

func (m *Monitor) PathStatusCached() {
	m.pathStatusCached.Inc()
}
