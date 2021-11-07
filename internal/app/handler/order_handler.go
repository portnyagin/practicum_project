package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -destination=mocks/mock_order_service.go -package=mocks . OrderService
type OrderService interface {
	Save(ctx context.Context, order *dto.Order) error
	GetOrderList(ctx context.Context, userID int) ([]dto.Order, error)
}

type OrderHandler struct {
	orderService OrderService
	auth         *Auth
	log          *infrastructure.Logger
}

func NewOrderHandler(os OrderService, auth *Auth, l *infrastructure.Logger) *OrderHandler {
	var target OrderHandler
	target.log = l
	target.orderService = os
	target.auth = auth
	return &target
}

/*
200 — номер заказа уже был загружен этим пользователем;
202 — новый номер заказа принят в обработку;
400 — неверный формат запроса;
401 — пользователь не аутентифицирован;
409 — номер заказа уже был загружен другим пользователем;
422 — неверный формат номера заказа;
500 — внутренняя ошибка сервера.
*/
func (h *OrderHandler) RegisterNewOrder(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		h.log.Error("OrderHandler:can't get request body", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}

	if len(b) == 0 || r.Header.Get("Content-Type") != "text/plain" {
		h.log.Info("OrderHandler:empty request body")
		if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("Неверный формат запроса")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}
	var (
		order dto.Order
	)

	order.Num = string(b)
	ctx := r.Context()
	u, _, err := h.auth.GetFromContext(ctx)
	if err != nil {
		h.log.Error("OrderHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}
	order.UserID = u
	err = h.orderService.Save(ctx, &order)
	if err != nil {
		h.log.Error("OrderHandler:recieved an error", zap.Error(err))
		if errors.Is(err, dto.ErrOrderRegistered) {
			if err = WriteResponse(w, http.StatusOK, ErrMessage("Номер заказа уже был загружен этим пользователем")); err != nil {
				h.log.Error("OrderHandler: can't write response", zap.Error(err))
			}
			return
		} else if errors.Is(err, dto.ErrOrderRegisteredByAnotherUser) {
			if err = WriteResponse(w, http.StatusConflict, ErrMessage("Номер заказа уже был загружен другим пользователем")); err != nil {
				h.log.Error("OrderHandler: can't write response", zap.Error(err))
			}
			return
		} else if errors.Is(err, dto.ErrBadParam) {
			if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("Неверный формат номера заказа")); err != nil {
				h.log.Error("OrderHandler: can't write response", zap.Error(err))
			}
			return
		} else if errors.Is(err, dto.ErrBadOrderNum) {
			if err = WriteResponse(w, http.StatusUnprocessableEntity, ErrMessage("Неверный формат номера заказа")); err != nil {
				h.log.Error("OrderHandler: can't write response", zap.Error(err))
			}
			return
		} else {
			if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
				h.log.Error("OrderHandler: can't write response", zap.Error(err))
			}
			return
		}
	}
	if err = WriteResponse(w, http.StatusAccepted, ""); err != nil {
		h.log.Error("OrderHandler: can't write response", zap.Error(err))
	}
	h.log.Info(fmt.Sprintf("Order %s succefully registered", order.Num))
}

/*
200 — успешная обработка запроса.
204 — нет данных для ответа.
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера
*/
func (h *OrderHandler) GetOrderList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _, err := h.auth.GetFromContext(ctx)
	if err != nil {
		h.log.Error("OrderHandler:can't get params from the token", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}

	res, err := h.orderService.GetOrderList(ctx, userID)
	if err != nil {
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
	} else if len(res) == 0 {
		if err = WriteResponse(w, http.StatusNoContent, ErrMessage("Нет данных для ответа")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
	} else {
		responseBody, err := json.Marshal(res)
		if err != nil {
			h.log.Error("OrderHandler: can't serialize response", zap.Error(err))
			return
		}
		if err = WriteResponse(w, http.StatusOK, responseBody); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
	}
}
