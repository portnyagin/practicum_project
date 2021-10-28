package handler

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/handler/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandler_Register(t *testing.T) {
	type wants struct {
		responseCode int
		contentType  string
	}
	type args struct {
		body      string
		login     string
		pass      string
		wantError bool
		err       error
	}
	tests := []struct {
		name  string
		wants wants
		args  args
	}{
		{name: "AuthHandler. Register. Case #1/ Positive",
			wants: wants{
				responseCode: http.StatusOK,
				contentType:  "application/json",
			},
			args: args{
				body:  "{\"login\": \"%s\",\"password\": \"%s\"}",
				login: "userLogin",
				pass:  "userPass",
				err:   nil,
			},
		},
		{name: "AuthHandler. Register. Case #2. Bad Query",
			wants: wants{
				responseCode: http.StatusBadRequest,
				contentType:  "application/json",
			},
			args: args{
				body:  "{\"model\": \"model x\",\"brand\": \"tesla\"}",
				login: "",
				pass:  "",
				err:   dto.ErrBadParam,
			},
		},
		{name: "AuthHandler. Register. Case #3. Empty Query",
			wants: wants{
				responseCode: http.StatusBadRequest,
				contentType:  "application/json",
			},
			args: args{
				body:  "",
				login: "",
				pass:  "",
				err:   dto.ErrBadParam,
			},
		},
		{name: "AuthHandler. Register. Case #3. Login Busy",
			wants: wants{
				responseCode: http.StatusConflict,
				contentType:  "application/json",
			},
			args: args{
				body:  "{\"login\": \"%s\",\"password\": \"%s\"}",
				login: "duplicateLogin",
				pass:  "anyPass",
				err:   dto.ErrDuplicateKey,
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	authService := mocks.NewMockAuthService(mockCtrl)
	target := NewAuthHandler(authService, log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService.EXPECT().Register(&dto.User{Login: tt.args.login, Pass: tt.args.pass}).Return(tt.args.err)
			body := strings.NewReader(fmt.Sprintf(tt.args.body, tt.args.login, tt.args.pass))
			request := httptest.NewRequest("POST", "/api/user/register", body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.Register)
			h.ServeHTTP(w, request)
			res := w.Result()
			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected status %d, got %d", tt.wants.contentType, contentType)
		})
	}
}
