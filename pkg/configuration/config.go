package configuration

import (
	"log/syslog"

	"github.com/caarlos0/env"
)

// ServiceConfig service configuration
type ServiceConfig struct {
	Server   ServerConfig
	Security SecurityConfig
	Redis    RedisConfig
}

// ServerConfig server configuration
type ServerConfig struct {
	HTTP     string          `env:"HTTP_ADDR" envDefault:":8077"`
	GRPC     string          `env:"HTTP_ADDR" envDefault:":8088"`
	LogLevel syslog.Priority `env:"LOG_LEVEL" envDefault:"7"` // LOG_DEBUG // LOG_ERR = 3
}

// SecurityConfig security configuration
type SecurityConfig struct {
	JWTSecret                    string `env:"JWT_SECRET_KEY"`
	JWTTokenExpiration           int    `env:"JWT_EXPIRE_HOURS" envDefault:"10"`
	JWTRefreshExpiration         int    `env:"JWT_REFRESH_HOURS" envDefault:"20"`
	JWTInternalSecret            string `env:"JWT_INTERNAL_SECRET_KEY"`
	JWTInternalTokenExpiration   int    `env:"JWT_INTERNAL_EXPIRE_HOURS" envDefault:"2"`
	JWTInternalRefreshExpiration int    `env:"JWT_INTERNAL_REFRESH_HOURS" envDefault:"5"`
}

// RedisConfig redis configuration
type RedisConfig struct {
	Addr string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	Pass string `env:"REDIS_PASS" envDefault:""`
	DB   int    `env:"REDIS_DB" evnDefault:"10"`
}

// Read service configuration from environment varible
func Read() (*ServiceConfig, error) {
	config := ServiceConfig{}

	if err := env.Parse(&config.Server); err != nil {
		return nil, err
	}
	if err := env.Parse(&config.Security); err != nil {
		return nil, err
	}
	if err := env.Parse(&config.Redis); err != nil {
		return nil, err
	}
	return &config, nil
}
