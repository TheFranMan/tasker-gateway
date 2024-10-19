package cache

import (
	"github.com/TheFranMan/tasker-common/types"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) GetKey(key string) (*types.RequestStatusString, error) {
	args := m.Called(key)
	return args.Get(0).(*types.RequestStatusString), args.Error(1)
}

func (m *Mock) SetKey(key string, value types.RequestStatusString) error {
	args := m.Called(key, value)
	return args.Error(0)
}
