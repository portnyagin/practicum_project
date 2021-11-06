package handler

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/handler/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestOrderHandler_RegisterNewOrder(t *testing.T) {
	type wants struct {
		responseCode int
		contentType  string
	}
	type args struct {
		body      string
		wantError bool
		err       error
	}
	tests := []struct {
		name  string
		wants wants
		args  args
	}{
		{name: "OrderHandler. RegisterNewOrder. Case #1. Positive",
			wants: wants{
				responseCode: http.StatusAccepted,
				contentType:  "application/json",
			},
			args: args{
				body:      "123456",
				wantError: false,
				err:       nil,
			},
		},
		{name: "OrderHandler. RegisterNewOrder. Case #2. Empty body",
			wants: wants{
				responseCode: http.StatusBadRequest,
				contentType:  "application/json",
			},
			args: args{
				body:      "",
				wantError: false,
				err:       dto.ErrBadParam,
			},
		},
		{name: "OrderHandler. RegisterNewOrder. Case #3. Invalid order num",
			wants: wants{
				responseCode: http.StatusUnprocessableEntity,
				contentType:  "application/json",
			},
			args: args{
				body: "Invalid order num",
				err:  dto.ErrBadParam,
			},
		},
		{name: "OrderHandler. RegisterNewOrder. Case #4. Order registered early by current user",
			wants: wants{
				responseCode: http.StatusOK,
				contentType:  "application/json",
			},
			args: args{
				body: "123456789",
				err:  dto.ErrOrderRegistered,
			},
		},
		{name: "OrderHandler. RegisterNewOrder. Case #5. Order registered early by another user",
			wants: wants{
				responseCode: http.StatusConflict,
				contentType:  "application/json",
			},
			args: args{
				body: "123456789",
				err:  dto.ErrOrderRegisteredByAnotherUser,
			},
		},
		{name: "OrderHandler. RegisterNewOrder. Case #6. Internal error",
			wants: wants{
				responseCode: http.StatusInternalServerError,
				contentType:  "application/json",
			},
			args: args{
				body: "123456789",
				err:  errors.New("any error"),
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	orderService := mocks.NewMockOrderService(mockCtrl)
	target := NewOrderHandler(orderService, auth, log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			orderService.EXPECT().Save(ctx, &dto.Order{Num: tt.args.body, UserID: 0}).Return(tt.args.err)
			body := strings.NewReader(tt.args.body)
			request := httptest.NewRequest("POST", "/api/user/orders", body)
			request.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.RegisterNewOrder)
			h.ServeHTTP(w, request)
			res := w.Result()
			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected status %d, got %d", tt.wants.contentType, contentType)
		})
	}
}

func generateOrderList(cnt int, userID int) []dto.Order {
	var res []dto.Order
	for i := 0; i < cnt; i++ {
		var obj dto.Order
		obj.UserID = userID
		obj.Num = strconv.Itoa(i)
		obj.Status = "NEW"
		obj.Accrual = 12
		obj.UploadAt = time.Now()

		res = append(res, obj)
	}
	return res
}

func TestOrderHandler_GetOrderList(t *testing.T) {
	type wants struct {
		responseCode int
		contentType  string
	}
	type args struct {
		userID    int
		objCount  int
		wantError bool
		err       error
	}
	tests := []struct {
		name  string
		wants wants
		args  args
	}{
		{name: "OrderHandler. RegisterNewOrder. Case #1. Positive",
			wants: wants{
				responseCode: http.StatusOK,
				contentType:  "application/json",
			},
			args: args{
				wantError: false,
				objCount:  5,
				err:       nil,
			},
		},
		{name: "OrderHandler. RegisterNewOrder. Case #2. No data for response",
			wants: wants{
				responseCode: http.StatusNoContent,
				contentType:  "application/json",
			},
			args: args{
				wantError: false,
				objCount:  0,
				err:       nil,
			},
		},
		{name: "OrderHandler. RegisterNewOrder. Case #3. Internal error",
			wants: wants{
				responseCode: http.StatusInternalServerError,
				contentType:  "application/json",
			},
			args: args{
				wantError: false,
				objCount:  0,
				err:       errors.New("any error"),
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	orderService := mocks.NewMockOrderService(mockCtrl)

	target := NewOrderHandler(orderService, auth, log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			orderService.EXPECT().
				GetOrderList(ctx, tt.args.userID).
				Return(generateOrderList(tt.args.objCount, tt.args.userID), tt.args.err)

			request := httptest.NewRequest("GET", "/api/user/orders", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.GetOrderList)
			h.ServeHTTP(w, request)
			res := w.Result()

			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected status %d, got %d", tt.wants.contentType, contentType)
		})
	}
}
