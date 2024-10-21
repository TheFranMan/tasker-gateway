package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setEnvs() {
	os.Setenv("PORT", "3000")
	os.Setenv("DB_USER", "DB_USER")
	os.Setenv("DB_PASS", "DB_PASS")
	os.Setenv("DB_HOST", "DB_HOST")
	os.Setenv("DB_PORT", "1")
	os.Setenv("DB_NAME", "DB_NAME")
	os.Setenv("AUTH_TOKENS", "AUTH_TOKENS")
	os.Setenv("REDIS_ADDR", "REDIS_ADDR")
}

func unsetEnvs() {
	os.Unsetenv("PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("AUTH_TOKENS")
	os.Unsetenv("REDIS_ADDR")
}

func TestConfig(t *testing.T) {
	setEnvs()
	defer unsetEnvs()

	t.Run("Can set env locale", func(t *testing.T) {
		for _, test := range []struct {
			env  string
			want Envs
		}{
			{
				env: "Prod",
				want: Envs{
					IsLocal: false,
					IsStage: false,
					IsProd:  true,
				},
			},
			{
				env: "prod",
				want: Envs{
					IsLocal: false,
					IsStage: false,
					IsProd:  true,
				},
			},
			{
				env: "Production",
				want: Envs{
					IsLocal: false,
					IsStage: false,
					IsProd:  true,
				},
			},
			{
				env: "Stage",
				want: Envs{
					IsLocal: false,
					IsStage: true,
					IsProd:  false,
				},
			},
			{
				env: "stage",
				want: Envs{
					IsLocal: false,
					IsStage: true,
					IsProd:  false,
				},
			},
			{
				env: "Staging",
				want: Envs{
					IsLocal: false,
					IsStage: true,
					IsProd:  false,
				},
			},
			{
				env: "Local",
				want: Envs{
					IsLocal: true,
					IsStage: false,
					IsProd:  false,
				},
			},
			{
				env: "Dev",
				want: Envs{
					IsLocal: true,
					IsStage: false,
					IsProd:  false,
				},
			},
			{
				env: "",
				want: Envs{
					IsLocal: true,
					IsStage: false,
					IsProd:  false,
				},
			},
		} {
			t.Run(test.env, func(t *testing.T) {
				os.Setenv("ENV", test.env)

				cfg, err := GetConfig()
				require.Nil(t, err)

				require.Equal(t, test.want.IsLocal, cfg.IsLocal)
				require.Equal(t, test.want.IsStage, cfg.IsStage)
				require.Equal(t, test.want.IsProd, cfg.IsProd)

				os.Unsetenv("ENV")
			})
		}

	})
}
