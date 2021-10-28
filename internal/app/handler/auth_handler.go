package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -destination=mocks/mock_auth_service.go -package=mocks . AuthService
type AuthService interface {
	Register(user *dto.User) error
}

type AuthHandler struct {
	authService AuthService
	log         *infrastructure.Logger
}

func NewAuthHandler(as AuthService, l *infrastructure.Logger) *AuthHandler {
	var target AuthHandler
	target.log = l
	target.authService = as

	return &target
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	b, err := getRequestBody(r)
	if err != nil {
		h.log.Error("AuthHandler:can't get request body", zap.Error(err))
		if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	if err := json.Unmarshal(b, &user); err != nil {
		h.log.Error("AuthHandler:can't unmarshal body", zap.Error(err))
		if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("Неверный формат запроса")); err != nil {
			h.log.Error("AuthHandler: can't write response", zap.Error(err))
		}
		return
	}
	err = h.authService.Register(&user)
	if err != nil {
		h.log.Error("AuthHandler:recieved an error", zap.Error(err))
		if errors.Is(err, dto.ErrDuplicateKey) {
			if err = WriteResponse(w, http.StatusConflict, ErrMessage("Логин уже занят")); err != nil {
				h.log.Error("AuthHandler: can't write response", zap.Error(err))
			}
			return
		} else if errors.Is(err, dto.ErrBadParam) {
			if err = WriteResponse(w, http.StatusBadRequest, ErrMessage("Неверный формат запроса")); err != nil {
				h.log.Error("AuthHandler: can't write response", zap.Error(err))
			}
			return
		} else {
			if err = WriteResponse(w, http.StatusInternalServerError, ErrMessage("Внутренняя ошибка сервера")); err != nil {
				h.log.Error("AuthHandler: can't write response", zap.Error(err))
			}
			return
		}
	}
	if err = WriteResponse(w, http.StatusOK, ""); err != nil {
		h.log.Error("AuthHandler: can't write response", zap.Error(err))
	}
	h.log.Info(fmt.Sprintf("User %s succefully registered", user.Login))

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}
