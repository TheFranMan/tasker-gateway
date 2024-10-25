package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) PathStatusCode(path string, code int) {
	m.Called(path, code)
}

func (m *Mock) StatusCacheHit() {
	m.Called()
}

func (m *Mock) StatusCacheMiss() {
	m.Called()
}

func (m *Mock) StatusDurationStart() *prometheus.Timer {
	m.Called()
	return &prometheus.Timer{}
}

func (m *Mock) StatusDurationEnd(timer *prometheus.Timer) {
	m.Called(timer)
}
