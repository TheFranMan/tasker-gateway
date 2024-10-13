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

	t.Run("token", func(t *testing.T) {
		for _, test := range []struct {
			want  bool
			token string
		}{
			{false, ""},
			{false, "nope"},
			{false, "acea7232-8a15-44a8-9808-7bf178b3cc5"},   // Too short
			{false, "acea7232-8a15-44a8-9808-7bf178b3cc566"}, // Too long
			{false, "acea72328a1544a898087bf178b3cc56"},

			{true, "acea7232-8a15-44a8-9808-7bf178b3cc56"},
		} {
			t.Run(test.token, func(t *testing.T) {
				require.Equal(t, test.want, ValidToken(test.token))
			})
		}
	})
}
