package repo

import (
	"fmt"
	"testing"

	"github.com/TheFranMan/tasker-common/types"
	"github.com/stretchr/testify/require"
)

func Test_can_convert_request_status_to_string(t *testing.T) {
	for _, test := range []struct {
		status types.RequestStatus
		want   types.RequestStatusString
	}{
		{types.RequestStatusNew, types.RequestStatusStringNew},
		{types.RequestStatusInProgress, types.RequestStatusStringInProgress},
		{types.RequestStatusFailed, types.RequestStatusStringFailed},
		{types.RequestStatusCompleted, types.RequestStatusStringCompleted},
	} {
		t.Run(fmt.Sprintf("%v", test.status), func(t *testing.T) {
			require.Equal(t, test.want, getRequestStatusString(test.status))
		})
	}
}
