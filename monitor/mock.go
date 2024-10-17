package monitor

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

func (m *Mock) PathStatusCode(path string, code int) {
	m.Called(path, code)
}

func (m *Mock) PathStatusCached() {
	m.Called()
}
