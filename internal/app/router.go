package app

import (
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/portnyagin/practicum_project/internal/app/handler"
)

func publicRoutes(r chi.Router, handler *handler.AuthHandler) {
	r.Group(func(router chi.Router) {
		router.Use(middleware.CleanPath)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Post("/api/user/register", handler.Register)
		router.Post("/api/user/login", handler.Login)
	})
}

func protectedOrderRoutes(r chi.Router, tokenAuth *jwtauth.JWTAuth, handler *handler.OrderHandler) {
	r.Group(func(router chi.Router) {
		router.Use(middleware.CleanPath)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)
		router.Post("/api/user/orders", handler.RegisterNewOrder)
		router.Get("/api/user/orders", handler.GetOrderList)
	})
}

func protectedBalanceRoutes(r chi.Router, tokenAuth *jwtauth.JWTAuth, handler *handler.BalanceHandler) {
	r.Group(func(router chi.Router) {
		router.Use(middleware.CleanPath)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)
		router.Get("/api/user/balance", handler.GetBalance)
		router.Post("/api/user/balance/withdraw", handler.Withdraw)
		router.Get("/api/user/balance/withdrawals", handler.GetWithdrawalsList)
	})
}
