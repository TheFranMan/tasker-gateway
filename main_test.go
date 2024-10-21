package main

import (
	"gateway/common"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_server_port(t *testing.T) {
	t.Run("local enviroment", func(t *testing.T) {
		config := &common.Config{
			Port: 3000,
		}
		config.IsLocal = true

		require.Equal(t, "localhost:3000", getServerPort(config))
	})

	t.Run("non local enviroment", func(t *testing.T) {
		config := &common.Config{
			Port: 3000,
		}
		config.IsProd = true

		require.Equal(t, ":3000", getServerPort(config))
	})
}
