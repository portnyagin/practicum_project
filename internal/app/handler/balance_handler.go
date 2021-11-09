package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -destination=mocks/mock_balance_service.go -package=mocks . BalanceService
type BalanceService interface {
	GetCurrentBalance(ctx context.Context, userID int) (*dto.Balance, error)
	Withdraw(ctx context.Context, obj *dto.Withdraw, userID int) error
	GetWithdrawalsList(ctx context.Context, userID int) ([]dto.Withdrawal, error)
}

type BalanceHandler struct {
	balanceService BalanceService
	auth           *Auth
	log            *infrastructure.Logger
}

func NewBalanceHandler(bs BalanceService, auth *Auth, l *infrastructure.Logger) *BalanceHandler {
	var target BalanceHandler
	target.log = l
	target.balanceService = bs
	target.auth = auth
	return &target
}

/*
200 — успешная обработка запроса
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/
func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, err := h.auth.GetFromContext(ctx)
	if err != nil {
		h.log.Error("BalanceHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	balance, err := h.balanceService.GetCurrentBalance(ctx, userID)
	if err != nil {
		h.log.Error("BalanceHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	responseBody, err := json.Marshal(balance)
	if err != nil {
		h.log.Error("BalanceHandler: can't serialize response", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}
	if err = WriteResponse(w, http.StatusOK, responseBody); err != nil {
		h.log.Error("BalanceHandler: can't write response", zap.Error(err))
	}
}

/*
200 — успешная обработка запроса;
401 — пользователь не авторизован;
402 — на счету недостаточно средств;
422 — неверный номер заказа;
500 — внутренняя ошибка сервера.
*/
func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		h.log.Error("BalanceHandler:can't get request body", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}
	if len(b) == 0 || r.Header.Get("Content-Type") != "application/json" {
		h.log.Info("BalanceHandler:empty request body")
		if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("Неверный формат запроса")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}

	var (
		withdraw dto.Withdraw
	)
	err = json.Unmarshal(b, &withdraw)
	if err != nil {
		h.log.Error("BalanceHandler:can't unmarshal request body", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}
	ctx := r.Context()
	userID, _, err := h.auth.GetFromContext(ctx)
	h.log.Info("Try to withdraw funds", zap.String("OrderNum", withdraw.OrderNum), zap.Int("userID", userID))
	if err != nil {
		h.log.Error("BalanceHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
		return
	}
	err = h.balanceService.Withdraw(ctx, &withdraw, userID)
	if err != nil {
		h.log.Error("BalanceHandler:Withdraw error", zap.Error(err))
		if errors.Is(err, dto.ErrNotEnoughFunds) {
			if err = WriteResponse(w, http.StatusPaymentRequired, ErrMessage("На счету недостаточно средств")); err != nil {
				h.log.Error("BalanceHandler: can't write response", zap.Error(err))
			}
		} else if errors.Is(err, dto.ErrBadOrderNum) {
			if err = WriteResponse(w, http.StatusUnprocessableEntity, ErrMessage("Неверный номер заказа")); err != nil {
				h.log.Error("BalanceHandler: can't write response", zap.Error(err))
			}
		} else {
			if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
				h.log.Error("BalanceHandler: can't write response", zap.Error(err))
			}
		}
		return
	}
	if err = WriteResponse(w, http.StatusOK, nil); err != nil {
		h.log.Error("BalanceHandler: can't write response", zap.Error(err))
	}
	h.log.Info("Withdraw success", zap.String("OrderNum", withdraw.OrderNum), zap.Int("userID", userID))
}

/*
200 — успешная обработка запроса
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/
func (h *BalanceHandler) GetWithdrawalsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, err := h.auth.GetFromContext(ctx)
	if err != nil {
		h.log.Error("BalanceHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	res, err := h.balanceService.GetWithdrawalsList(ctx, userID)
	if err != nil {
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
	} else {
		responseBody, err := json.Marshal(res)
		if err != nil {
			h.log.Error("BalanceHandler: can't serialize response", zap.Error(err))
			if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
				h.log.Error("BalanceHandler: can't write response", zap.Error(err))
			}
			return
		}
		if err = WriteResponse(w, http.StatusOK, responseBody); err != nil {
			h.log.Error("BalanceHandler: can't write response", zap.Error(err))
		}
	}
}
