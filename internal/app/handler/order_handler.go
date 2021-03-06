package handler

import (
	"context"
	"encoding/json"
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

	var (
		statusCode int
		msg        string
	)
	if err != nil {
		h.log.Error("OrderHandler:recieved an error", zap.Error(err))
		switch err {
		case dto.ErrOrderRegistered:
			statusCode = http.StatusOK
			msg = "Номер заказа уже был загружен этим пользователем"
		case dto.ErrOrderRegisteredByAnotherUser:
			statusCode = http.StatusConflict
			msg = "Номер заказа уже был загружен этим пользователем"
		case dto.ErrBadParam:
			statusCode = http.StatusBadRequest
			msg = "Неверный формат запроса"
		case dto.ErrBadOrderNum:
			statusCode = http.StatusUnprocessableEntity
			msg = "Неверный формат номера заказа"
		default:
			statusCode = http.StatusInternalServerError
			msg = "Внутренняя ошибка сервера"
		}
		if err = WriteResponse(w, statusCode, ErrMessage(msg)); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}

	if err = WriteResponse(w, http.StatusAccepted, nil); err != nil {
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
		h.log.Debug("/api/user/orders result", zap.String("dto", string(responseBody)))
		if err != nil {
			h.log.Error("OrderHandler: can't serialize response", zap.Error(err))
			return
		}
		if err = WriteResponse(w, http.StatusOK, responseBody); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
	}
}
