package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validation(t *testing.T) {
	t.Run("id", func(t *testing.T) {
		for _, test := range []struct {
			want bool
			id   int
		}{
			{want: false, id: -1},
			{want: false, id: 0},

			{want: true, id: 1},
			{want: true, id: 123456},
		} {
			t.Run(fmt.Sprintf("%d", test.id), func(t *testing.T) {
				require.Equal(t, test.want, ValidID(test.id))
			})
		}
	})
}
