package config

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	Database Database `mapstructure:",squash"`
	Log      Log      `mapstructure:",squash"`
	Server   Server   `mapstructure:",squash"`
}

// Database holds configurations to be used when connecting to a PostgreSQL server.
// URL can have the following parameters regarding pool configuration:
//   - pool_health_check_period: duration string
//   - pool_max_conn_idle_time: duration string
//   - pool_max_conn_lifetime: duration string
//   - pool_max_conns: integer greater than 0
//   - pool_min_conns: integer 0 or greater
type Database struct {
	DSN string `mapstructure:"database_dsn"`
}

type Log struct {
	Level string `mapstructure:"log_level"`
}

type Server struct {
	GracefulShutdownPeriod time.Duration `mapstructure:"server_graceful_shutdown_period"`
	Port                   string        `mapstructure:"server_port"`
	ReadHeaderTimeout      time.Duration `mapstructure:"server_read_header_timeout"`
}

func New() (*Config, error) {
	viper.SetDefault("DATABASE_DSN", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable&pool_health_check_period=1h&pool_max_conn_idle_time=2m&pool_max_conn_lifetime=30m&pool_max_conns=50&pool_min_conns=5")

	viper.SetDefault("LOG_LEVEL", "debug")

	viper.SetDefault("SERVER_GRACEFUL_SHUTDOWN_PERIOD", "30s")
	viper.SetDefault("SERVER_PORT", "3000")
	viper.SetDefault("SERVER_READ_HEADER_TIMEOUT", "2s")

	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		println("config: .env file not found")
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		),
	)); err != nil {
		return nil, err
	}

	return cfg, nil
}
