package redis

import (
	"strconv"

	"github.com/go-redis/redis"
	"github.com/im-kulikov/helium/module"
	"github.com/spf13/viper"
)

type (
	// Config alias
	Config = redis.UniversalOptions

	// Client alias
	Client = redis.UniversalClient

	// Error is constant error
	Error string
)

const (
	// ErrEmptyConfig when given empty options
	ErrEmptyConfig = Error("missing redis config key")
	// ErrEmptyAddresses when given empty addresses
	ErrEmptyAddresses = Error("missing addresses")
	// ErrPemParse when couldn't parse pem in sslrootcert
	ErrPemParse = Error("couldn't parse pem in sslrootcert")
	// ErrEmptyLogger when logger not initialized
	ErrEmptyLogger = Error("database empty logger")
	// ErrSSLKeyHasWorldPermissions when pk permissions no u=rw (0600) or less
	ErrSSLKeyHasWorldPermissions = Error("private key file has group or world access. Permissions should be u=rw (0600) or less")

	errUnsupportedSSLMode = `unsupported sslmode %q; only "require" (default), "verify-full", "verify-ca", and "disable" supported`
)

var (
	// Module is default Redis client
	Module = module.Module{
		{Constructor: NewDefaultConfig},
		{Constructor: NewConnection},
	}
)

// Error represented in string
func (e Error) Error() string { return string(e) }

// NewDefaultConfig for connection
func NewDefaultConfig(v *viper.Viper) (*Config, error) {
	if !v.IsSet("redis.address") && !v.IsSet("redis.addresses") && !v.IsSet("redis.addresses_0") {
		return nil, ErrEmptyConfig
	}

	var addresses []string

	if addr := v.GetString("redis.address"); addr != "" {
		addresses = append(addresses, addr)
	} else if addresses = fetchAddresses(v); len(addresses) == 0 {
		return nil, ErrEmptyAddresses
	}

	v.SetDefault("redis.options.sslmode", "disable")

	// re-fetch by full key
	options := v.GetStringMapString("redis.options")
	if len(options) > 0 {
		for opt := range options {
			options[opt] = v.GetString("redis.options." + opt)
		}
	}

	tlsConfig, err := ssl(options)
	if err != nil {
		return nil, err
	}

	return &Config{
		Addrs:              addresses,
		DB:                 v.GetInt("redis.db"),
		Password:           v.GetString("redis.password"),
		MaxRetries:         v.GetInt("redis.max_retries"),
		MinRetryBackoff:    v.GetDuration("redis.min_retry_backoff"),
		MaxRetryBackoff:    v.GetDuration("redis.max_retry_backoff"),
		DialTimeout:        v.GetDuration("redis.dial_timeout"),
		ReadTimeout:        v.GetDuration("redis.read_timeout"),
		WriteTimeout:       v.GetDuration("redis.write_timeout"),
		PoolSize:           v.GetInt("redis.pool_size"),
		MinIdleConns:       v.GetInt("redis.min_idle_cons"),
		MaxConnAge:         v.GetDuration("redis.max_con_age"),
		PoolTimeout:        v.GetDuration("redis.pool_timeout"),
		IdleTimeout:        v.GetDuration("redis.idle_timeout"),
		IdleCheckFrequency: v.GetDuration("redis.idle_check_frequency"),
		MaxRedirects:       v.GetInt("redis.max_redirects"),
		ReadOnly:           v.GetBool("redis.read_only"),
		RouteByLatency:     v.GetBool("redis.router_by_latency"),
		RouteRandomly:      v.GetBool("redis.router_randomly"),
		MasterName:         v.GetString("redis.master_name"),
		TLSConfig:          tlsConfig,
	}, nil
}

func fetchAddresses(v *viper.Viper) []string {
	var (
		addresses []string
	)

	for i := 0; ; i++ {
		addr := v.GetString("redis.addresses_" + strconv.Itoa(i))
		if addr == "" {
			break
		}

		addresses = append(addresses, addr)
	}

	if len(addresses) == 0 {
		addresses = v.GetStringSlice("redis.addresses")
	}

	return addresses
}

// NewConnection of redis client
func NewConnection(opts *Config) (cache Client, err error) {
	cache = redis.NewUniversalClient(opts)

	if _, err = cache.Ping().Result(); err != nil {
		return nil, err
	}

	return cache, nil
}
