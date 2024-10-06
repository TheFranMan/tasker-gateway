package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validation(t *testing.T) {
	t.Run("id", func(t *testing.T) {
		for _, test := range []struct {
			want bool
			id   string
		}{
			{want: false, id: ""},
			{want: false, id: "nope"},
			{want: false, id: "-1"},
			{want: false, id: "0"},
			{want: false, id: "1e"},
			{want: false, id: "e1"},
			{want: false, id: "1!"},
			{want: false, id: "!1"},

			{want: true, id: "1"},
			{want: true, id: "123456"},
		} {
			t.Run(test.id, func(t *testing.T) {
				require.Equal(t, test.want, ValidID(test.id))
			})
		}
	})
}
