package handler

import (
	"context"
	"errors"
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
		{name: "AuthHandler. Register. Case #1. Positive",
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
		{name: "AuthHandler. Register. Case #3. Bad Query2",
			wants: wants{
				responseCode: http.StatusBadRequest,
				contentType:  "application/json",
			},
			args: args{
				body:  "{\"login3\": \"login\"}",
				login: "",
				pass:  "",
				err:   dto.ErrBadParam,
			},
		},
		{name: "AuthHandler. Register. Case #4. Empty Query",
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
		{name: "AuthHandler. Register. Case #5. Login Busy",
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
		{name: "AuthHandler. Register. Case #6. Service Error",
			wants: wants{
				responseCode: http.StatusBadRequest,
				contentType:  "application/json",
			},
			args: args{
				body:  "{\"login\": \"%s\",\"password\": \"%s\"}",
				login: "anyLogin",
				pass:  "anyPass",
				err:   dto.ErrBadParam,
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	authService := mocks.NewMockAuthService(mockCtrl)

	target := NewAuthHandler(authService, auth, log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			authService.EXPECT().
				Register(ctx, &dto.User{Login: tt.args.login, Pass: tt.args.pass}).
				Return(&dto.User{ID: 10, Login: tt.args.login, Pass: tt.args.pass}, tt.args.err).
				AnyTimes()
			body := strings.NewReader(fmt.Sprintf(tt.args.body, tt.args.login, tt.args.pass))
			request := httptest.NewRequest("POST", "/api/user/register", body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.Register)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected status %d, got %d", tt.wants.contentType, contentType)

			if res.StatusCode == http.StatusOK {
				cookies := res.Cookies()
				var token *http.Cookie
				for _, c := range cookies {
					if c.Name == "jwt" {
						token = c
						break
					}
				}
				assert.NotNil(t, token, "JWT token not set")
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	type wants struct {
		responseCode int
		contentType  string
	}
	type args struct {
		body    string
		login   string
		pass    string
		allowed bool
		err     error
	}
	tests := []struct {
		name  string
		wants wants
		args  args
	}{
		{name: "AuthHandler. Check. Case #1. Positive",
			wants: wants{
				responseCode: http.StatusOK,
				contentType:  "application/json",
			},
			args: args{
				body:    "{\"login\": \"%s\",\"password\": \"%s\"}",
				login:   "userLogin",
				pass:    "userPass",
				allowed: true,
				err:     nil,
			},
		},
		{name: "AuthHandler. Check. Case #2. Bad Query",
			wants: wants{
				responseCode: http.StatusBadRequest,
				contentType:  "application/json",
			},
			args: args{
				body:    "{\"model\": \"model x\",\"brand\": \"tesla\"}",
				login:   "",
				pass:    "",
				allowed: false,
				err:     dto.ErrBadParam,
			},
		},
		{name: "AuthHandler. Check. Case #3. Access Denied",
			wants: wants{
				responseCode: http.StatusUnauthorized,
				contentType:  "application/json",
			},
			args: args{
				body:    "{\"login\": \"%s\",\"password\": \"%s\"}",
				login:   "userLogin3",
				pass:    "userPass3",
				allowed: false,
				err:     nil,
			},
		},
		{name: "AuthHandler. Check. Case #4. Empty Query",
			wants: wants{
				responseCode: http.StatusBadRequest,
				contentType:  "application/json",
			},
			args: args{
				body:    "",
				login:   "",
				pass:    "",
				allowed: false,
				err:     dto.ErrBadParam,
			},
		},

		{name: "AuthHandler. Check. Case #5. Service Error",
			wants: wants{
				responseCode: http.StatusInternalServerError,
				contentType:  "application/json",
			},
			args: args{
				body:    "{\"login\": \"%s\",\"password\": \"%s\"}",
				login:   "anyLogin",
				pass:    "anyPass",
				allowed: false,
				err:     errors.New("any error"),
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	authService := mocks.NewMockAuthService(mockCtrl)
	target := NewAuthHandler(authService, auth, log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var resObj *dto.User
			if tt.args.allowed {
				resObj = &dto.User{ID: 10, Login: tt.args.login, Pass: tt.args.pass}
			}
			authService.EXPECT().
				Check(ctx, &dto.User{Login: tt.args.login, Pass: tt.args.pass}).
				Return(resObj, tt.args.err).
				AnyTimes()

			body := strings.NewReader(fmt.Sprintf(tt.args.body, tt.args.login, tt.args.pass))
			request := httptest.NewRequest("POST", "/api/user/login", body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.Login)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected contentType %d, got %d", tt.wants.contentType, contentType)

			if res.StatusCode == http.StatusOK {
				cookies := res.Cookies()
				var token *http.Cookie
				for _, c := range cookies {
					if c.Name == "jwt" {
						token = c
						break
					}
				}
				assert.NotNil(t, token, "JWT token not set")
			}
		})
	}
}

func TestAuthHandler_RegisterInt(t *testing.T) {
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
		{name: "AuthHandler. Register. Case #1. Positive",
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
	}

	target := NewAuthHandler(authService, auth, log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(fmt.Sprintf(tt.args.body, tt.args.login, tt.args.pass))
			request := httptest.NewRequest("POST", "/api/user/register", body)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(target.Register)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			contentType := res.Header.Get("Content-type")
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, contentType, "Expected status %d, got %d", tt.wants.contentType, contentType)
		})
	}
}
