package main

import (
	"sync"

	"github.com/rs/zerolog"

	"github.com/terdia/mvp/internal/service/productservice"
	"github.com/terdia/mvp/internal/service/userservice"
)

type (
	application struct {
		wg             *sync.WaitGroup
		config         *config
		logger         *zerolog.Logger
		userService    userservice.UserService
		productService productservice.ProductService
	}

	config struct {
		AppPort int `env:"APP_PORT" envDefault:"4000"`
		Debug   string
		Db      db
		Cors    struct {
			TrustedOrigins []string `env:"CORS_ALLOWED" envSeparator:","`
		}
	}

	db struct {
		Dsn          string `env:"POSTGRES_DB_DSN"`
		MaxOpenConns int    `env:"DB_MAX_OPEN_CONN" envDefault:"25"`
		MaxIdleConns int    `env:"DB_MAX_IDLE_CONN" envDefault:"25"`
		MaxIdleTime  string `env:"DB_MAX_IDLE_TIME" envDefault:"15m"`
	}
)
