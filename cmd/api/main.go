package main

import (
	"os"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"

	"github.com/terdia/mvp/internal/repository/repositorypermission"
	"github.com/terdia/mvp/internal/repository/repositoryproduct"
	"github.com/terdia/mvp/internal/repository/repositorytoken"
	"github.com/terdia/mvp/internal/repository/repositoryuser"
	"github.com/terdia/mvp/internal/service/auth"
	"github.com/terdia/mvp/internal/service/productservice"
	"github.com/terdia/mvp/internal/service/transaction"
	"github.com/terdia/mvp/internal/service/userservice"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse env")
	}

	postgresDb, err := OpenDb(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open connection to db")
	}

	defer postgresDb.Close() //nolint
	logger.Printf("database connection pool established")

	tokenService := auth.NewTokenService(repositorytoken.NewTokenRepository(postgresDb))

	newUserService := userservice.NewUserService(
		repositoryuser.NewUserRepository(postgresDb),
		tokenService,
		repositorypermission.NewPermissionRepository(postgresDb),
	)

	newProductService := productservice.NewProductService(
		repositoryproduct.NewProductRepository(postgresDb),
	)

	app := &application{
		wg:                 new(sync.WaitGroup),
		config:             &cfg,
		logger:             &logger,
		userService:        newUserService,
		productService:     newProductService,
		transactionService: transaction.NewTransactionService(newUserService, newProductService),
	}

	err = app.serve()
	if err != nil {
		logger.Fatal().Err(err).Msg("App serve failed")
	}
}
