package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	cfg "github.com/portnyagin/practicum_project/internal/app/config"
	"github.com/portnyagin/practicum_project/internal/app/handler"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure/client"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure/postgres"
	"github.com/portnyagin/practicum_project/internal/app/repository"
	"github.com/portnyagin/practicum_project/internal/app/service"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func Start() {
	logger, err := zap.NewDevelopment()
	//logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	config := cfg.NewConfig()
	err = config.Init()
	if err != nil {
		logger.Fatal("can't init configuration", zap.Error(err))
	}

	postgresHandler, err := postgres.NewPostgresqlHandler(context.Background(), config.DatabaseDSN)
	if err != nil {
		logger.Fatal("can't create postgres handler", zap.Error(err))
	}
	postgresHandlerTx, err := postgres.NewPostgresqlHandlerTX(context.Background(), config.DatabaseDSN, logger)
	if err != nil {
		logger.Fatal("can't create postgres handler", zap.Error(err))
	}

	if config.Reinit {
		err = repository.ClearDatabase(context.Background(), postgresHandler)
		if err != nil {
			logger.Fatal("can't clear database structure", zap.Error(err))
			return
		}
	}
	err = repository.InitDatabase(context.Background(), postgresHandlerTx)
	if err != nil {
		logger.Fatal("can't init database structure", zap.Error(err))
		return
	}

	userRepository, err := repository.NewUserRepository(postgresHandlerTx, logger)
	if err != nil {
		logger.Fatal("can't init user repopsitory", zap.Error(err))
		return
	}

	orderRepository, err := repository.NewOrderRepository(postgresHandlerTx, logger)
	if err != nil {
		logger.Fatal("can't init order repopsitory", zap.Error(err))
		return
	}
	balanceRepository, err := repository.NewBalanceRepository(postgresHandlerTx, logger)
	if err != nil {
		logger.Fatal("can't init balance repopsitory", zap.Error(err))
		return
	}

	authService := service.NewAuthService(userRepository, logger)
	orderService := service.NewOrderService(orderRepository, logger)
	balanceService := service.NewBalanceService(balanceRepository, logger)
	auth := handler.NewAuth("secret")
	authHandler := handler.NewAuthHandler(authService, auth, logger)
	orderHandler := handler.NewOrderHandler(orderService, auth, logger)
	balanceHandler := handler.NewBalanceHandler(balanceService, auth, logger)
	router := chi.NewRouter()

	accrualClient := client.NewAccrualClient(config.AccrualServiceAddress, logger)
	gophermartClient := client.NewGophermartClient(config.ServerAddress, logger)
	accrualService := service.NewAccrualService(orderRepository, balanceRepository, accrualClient, gophermartClient, logger)
	accrualHandler := handler.NewAccrualHandler(accrualService, logger)

	publicRoutes(router, authHandler, accrualHandler, postgresHandlerTx, logger)
	protectedOrderRoutes(router, auth.GetJWTAuth(), postgresHandlerTx, orderHandler, logger)
	protectedBalanceRoutes(router, auth.GetJWTAuth(), postgresHandlerTx, balanceHandler, logger)

	go accrualService.StartProcessJob(1)
	err = http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		fmt.Println("can't start service")
		fmt.Println(err)
	}
}
