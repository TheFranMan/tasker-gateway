package handlers

import (
	"fmt"
	"testing"

	"github.com/TheFranMan/tasker-common/types"
	"github.com/stretchr/testify/require"
)

func Test_can_get_status_as_a_string(t *testing.T) {
	for _, test := range []struct {
		status types.RequestStatus
		want   types.RequestStatusString
	}{
		{types.RequestStatusNew, types.RequestStatusStringNew},
		{types.RequestStatusInProgress, types.RequestStatusStringInProgress},
		{types.RequestStatusCompleted, types.RequestStatusStringCompleted},
		{types.RequestStatusFailed, types.RequestStatusStringFailed},
	} {
		t.Run(fmt.Sprintf("%d", test.status), func(t *testing.T) {
			require.Equal(t, test.want, getRequestStatusString(test.status))
		})
	}
}
