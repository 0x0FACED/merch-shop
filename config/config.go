package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type ServiceConfig struct {
	Server   ServerConfig
	Logger   LoggerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host string `env:"SERVER_HOST" envDefault:"localhost"`
	Port string `env:"SERVER_PORT" envDefault:"8080"`

	// http
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT"`

	// echo
	DebugMode  bool   `env:"SERVER_DEBUG_MODE"`
	CSRFSecret string `env:"SERVER_CSRF_TOKEN"`
}

type DatabaseConfig struct {
	DSN               string        `env:"DATABASE_DSN,required"`
	MaxOpenConns      int32         `env:"DATABASE_MAX_OPEN_CONNS"`
	MaxIdleConns      int32         `env:"DATABASE_MAX_IDLE_CONNS"`
	ConnMaxLifetime   time.Duration `env:"DATABASE_CONN_MAX_LIFETIME"`
	ConnMaxIdleTime   time.Duration `env:"DATABASE_CONN_MAX_IDLE_LIFETIME"`
	ConnectionTimeout time.Duration `env:"DATABASE_CONNECTION_TIMEOUT"`
	PoolTimeout       time.Duration `env:"DATABASE_POOL_TIMEOUT"`
}

type LoggerConfig struct {
	LogLevel string `env:"LOGGER_LEVEL" envDefault:"debug"`
}

// Паникует, если чет не получилось
func MustLoad() *ServiceConfig {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := &ServiceConfig{}

	if err := env.Parse(&cfg.Server); err != nil {
		panic(err)
	}

	if err := env.Parse(&cfg.Database); err != nil {
		panic(err)
	}

	if err := env.Parse(&cfg.Logger); err != nil {
		panic(err)
	}

	return cfg
}
