package redis

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func prepare() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}

func TestRedis(t *testing.T) {
	_ = ErrEmptyAddresses.Error()

	t.Run("Config", func(t *testing.T) {
		t.Run("should return error when config file is empty", func(t *testing.T) {
			cfg, err := NewDefaultConfig(viper.New())
			require.Empty(t, cfg)
			require.Error(t, err)
		})

		t.Run("should return config when address exists", func(t *testing.T) {
			v := prepare()
			v.SetDefault("redis.address", "localhost:6379")
			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)
			require.NotEmpty(t, cfg)
		})

		t.Run("should return config when address exists", func(t *testing.T) {
			v := prepare()
			v.SetDefault("redis.addresses_0", "localhost:6379")
			cfg, err := NewDefaultConfig(v)
			require.NoError(t, err)
			require.NotEmpty(t, cfg)
		})

		t.Run("should fail on empty address 1", func(t *testing.T) {
			v := prepare()
			v.SetDefault("redis.addresses", "")
			cfg, err := NewDefaultConfig(v)
			require.Error(t, err)
			require.Empty(t, cfg)
		})

		t.Run("should fail on empty address 2", func(t *testing.T) {
			v := prepare()
			v.SetDefault("redis.addresses_0", "")
			cfg, err := NewDefaultConfig(v)
			require.Error(t, err)
			require.Empty(t, cfg)
		})
	})

	t.Run("Connection", func(t *testing.T) {
		v := prepare()
		v.SetDefault("redis.address", "localhost:6379")
		cfg, err := NewDefaultConfig(v)
		require.NoError(t, err)
		require.NotEmpty(t, cfg)

		t.Run("should create redis client", func(t *testing.T) {
			cli, err := NewConnection(cfg)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.Ping()
			//
		})

		t.Run("should return error when address incorrect", func(t *testing.T) {
			cfg.Addrs = []string{"foo"}
			cli, err := NewConnection(cfg)
			require.Nil(t, cli)
			require.Error(t, err)
		})
	})
}
