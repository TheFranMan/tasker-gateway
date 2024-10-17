package cache

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

func (m *Mock) GetKey(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *Mock) SetKey(name string, value interface{}) error {
	args := m.Called(name, value)
	return args.Error(0)
}
