package handler

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/handler/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBalanceHandler_GetBalance(t *testing.T) {
	type args struct {
		balance *dto.Balance
		error   error
	}
	type wants struct {
		responseCode int
		contentType  string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "BalanceHandler. GetBalance. Case #1. Positive",
			args: args{
				balance: &dto.Balance{},
				error:   nil,
			},
			wants: wants{
				responseCode: http.StatusOK,
				contentType:  "application/json",
			},
		},
		{
			name: "BalanceHandler. GetBalance. Case #2. Nil result",
			args: args{
				balance: nil,
				error:   nil,
			},
			wants: wants{
				responseCode: http.StatusOK,
				contentType:  "application/json",
			},
		},
		{
			name: "BalanceHandler. GetBalance. Case #3. Error",
			args: args{
				balance: nil,
				error:   errors.New("any error"),
			},
			wants: wants{
				responseCode: http.StatusInternalServerError,
				contentType:  "application/json",
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	balanceService := mocks.NewMockBalanceService(mockCtrl)
	target := NewBalanceHandler(balanceService, auth, log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balanceService.EXPECT().GetCurrentBalance(gomock.Any(), 0).Return(tt.args.balance, tt.args.error)

			request := httptest.NewRequest("GET", "/api/user/balance", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.GetBalance)
			h.ServeHTTP(w, request)
			res := w.Result()

			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected status %d, got %d", tt.wants.contentType, contentType)
		})
	}
}

func TestBalanceHandler_Withdraw(t *testing.T) {
	type args struct {
		error error
		body  string
	}
	type wants struct {
		responseCode int
		contentType  string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "BalanceHandler. WithdrawalsList. Case #1. Positive",
			args: args{
				error: nil,
				body:  "{\"order\": \"2377225624\",\"sum\": 751}",
			},
			wants: wants{
				responseCode: http.StatusOK,
				contentType:  "application/json",
			},
		},
		{
			name: "BalanceHandler. WithdrawalsList. Case #2. NotEnoughFunds",
			args: args{
				error: dto.ErrNotEnoughFunds,
				body:  "{\"order\": \"2377225624\",\"sum\": 751}",
			},
			wants: wants{
				responseCode: http.StatusPaymentRequired,
				contentType:  "application/json",
			},
		},
		{
			name: "BalanceHandler. WithdrawalsList. Case #3. Any error",
			args: args{
				error: errors.New("Any error"),
				body:  "{\"order\": \"2377225624\",\"sum\": 751}",
			},
			wants: wants{
				responseCode: http.StatusInternalServerError,
				contentType:  "application/json",
			},
		},
		{
			name: "BalanceHandler. WithdrawalsList. Case #4. Bad Order Num",
			args: args{
				error: dto.ErrBadOrderNum,
				body:  "{\"order\": \"2377225624\",\"sum\": 751}",
			},
			wants: wants{
				responseCode: http.StatusUnprocessableEntity,
				contentType:  "application/json",
			},
		},
		{
			name: "BalanceHandler. WithdrawalsList. Case #5 Bad request (empty body)",
			args: args{
				error: dto.ErrBadOrderNum,
				body:  "",
			},
			wants: wants{
				responseCode: http.StatusBadRequest,
				contentType:  "application/json",
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	balanceService := mocks.NewMockBalanceService(mockCtrl)
	target := NewBalanceHandler(balanceService, auth, log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balanceService.EXPECT().Withdraw(gomock.Any(), gomock.Any(), 0).Return(tt.args.error)
			body := strings.NewReader(tt.args.body)
			request := httptest.NewRequest("POST", "/api/user/balance/withdraw", body)
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.Withdraw)
			h.ServeHTTP(w, request)
			res := w.Result()

			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected status %d, got %d", tt.wants.contentType, contentType)
		})
	}
}

func TestBalanceHandler_WithdrawalsList(t *testing.T) {
	type args struct {
		res   []dto.Withdrawal
		error error
	}
	type wants struct {
		responseCode int
		contentType  string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "BalanceHandler. WithdrawalsList. Case #1. Positive",
			args: args{
				res:   []dto.Withdrawal{dto.Withdrawal{}},
				error: nil,
			},
			wants: wants{
				responseCode: http.StatusOK,
				contentType:  "application/json",
			},
		},
		{
			name: "BalanceHandler. WithdrawalsList. Case #2. Error",
			args: args{
				res:   nil,
				error: errors.New("any error"),
			},
			wants: wants{
				responseCode: http.StatusInternalServerError,
				contentType:  "application/json",
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	balanceService := mocks.NewMockBalanceService(mockCtrl)
	target := NewBalanceHandler(balanceService, auth, log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balanceService.EXPECT().GetWithdrawalsList(gomock.Any(), 0).Return(tt.args.res, tt.args.error)

			request := httptest.NewRequest("GET", "/api/user/balance/withdrawals", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.GetWithdrawalsList)
			h.ServeHTTP(w, request)
			res := w.Result()

			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected status %d, got %d", tt.wants.contentType, contentType)
		})
	}
}
