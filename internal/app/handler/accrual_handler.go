package handler

import (
	"context"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

//go:generate mockgen -destination=mocks/mock_accrual_service.go -package=mocks . AccrualService
type AccrualService interface {
	ProcessOrder(ctx context.Context, orderNum string) error
}

type AccrualHandler struct {
	accrualService AccrualService
	log            *infrastructure.Logger
}

func NewAccrualHandler(accrualService AccrualService, l *infrastructure.Logger) *AccrualHandler {
	var target AccrualHandler
	target.accrualService = accrualService
	target.log = l
	return &target
}

func (h *AccrualHandler) ProcessOrder(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "" || r.RequestURI[1:] == "" {
		h.log.Error("AccrualHandler:bad request", zap.String("RequestURI", r.RequestURI))
		if err := WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))
		}
		return
	}
	url := strings.Split(r.RequestURI[1:], "/")
	orderNum := url[len(url)-1]
	ctx := r.Context()
	err := h.accrualService.ProcessOrder(ctx, orderNum)
	if err != nil {
		h.log.Error("AccrualHandler:internal service error", zap.String("RequestURI", r.RequestURI), zap.Error(err))
		if err := WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("OrderHandler: can't write response", zap.Error(err))

		}
		return
	}
	if err = WriteResponse(w, http.StatusCreated, nil); err != nil {
		h.log.Error("AccrualHandler: can't write response", zap.Error(err))
	}
}
