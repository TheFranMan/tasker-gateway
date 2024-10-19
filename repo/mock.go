package repo

import (
	"github.com/TheFranMan/tasker-common/types"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) NewDelete(authToken string, id int) (string, error) {
	args := m.Called(authToken, id)
	return args.String(0), args.Error(1)
}

func (m *Mock) GetStatus(token string) (*types.RequestStatusString, error) {
	args := m.Called(token)

	if nil == args.Get(0) {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.RequestStatusString), args.Error(1)
}
