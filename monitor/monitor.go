package monitor

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Interface interface {
	PathStatusCode(path string, code int)
}

type Monitor struct {
	pathStatusCodes *prometheus.CounterVec
}

func New() *Monitor {
	pathStatusCodes := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_path_status_codes",
		Help: "The http statuc code per path",
	}, []string{"path", "status_code"})

	prometheus.Register(pathStatusCodes)

	return &Monitor{
		pathStatusCodes: pathStatusCodes,
	}
}

func (m *Monitor) PathStatusCode(path string, code int) {
	m.pathStatusCodes.WithLabelValues(path, strconv.Itoa(code)).Inc()
}
