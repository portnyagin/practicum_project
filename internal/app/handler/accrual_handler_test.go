package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/portnyagin/practicum_project/internal/app/handler/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccrualHandler_ProcessOrder(t *testing.T) {
	type wants struct {
		responseCode int
		contentType  string
	}
	type args struct {
		orderNum string
	}
	tests := []struct {
		name  string
		wants wants
		args  args
	}{
		{name: "AuthHandler. Register. Case #1. Positive",
			wants: wants{
				responseCode: http.StatusOK,
				contentType:  "application/json",
			},
			args: args{
				orderNum: "1223",
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	accrualService := mocks.NewMockAccrualService(mockCtrl)

	target := NewAccrualHandler(accrualService, log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*ctx := context.Background()
			accrualService.EXPECT().
				Register(ctx, &dto.User{Login: tt.args.login, Pass: tt.args.pass}).
				Return(&dto.User{ID: 10, Login: tt.args.login, Pass: tt.args.pass}, tt.args.err).
				AnyTimes()*/
			request := httptest.NewRequest("POST", "/api/orders/"+tt.args.orderNum, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.ProcessOrder)
			h.ServeHTTP(w, request)
			res := w.Result()
			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected status %d, got %d", tt.wants.contentType, contentType)

		})
	}
}
