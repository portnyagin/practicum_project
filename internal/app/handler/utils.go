package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/jwtauth/v5"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"io"
	"net/http"
)

func ErrMessage(msg string) dto.Error {
	return dto.Error{Msg: msg}
}

func WriteResponse(w http.ResponseWriter, status int, message interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		panic("Can't write response")
	}
	return nil
}

func getRequestBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Auth struct {
	tokenAuth *jwtauth.JWTAuth
}

func NewAuth(secret string) *Auth {
	var auth Auth
	auth.tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
	return &auth
}

func (auth *Auth) GetFromContext(ctx context.Context) (userID int, login string, err error) {
	_, m, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, "", err
	}

	if u, ok := m["user_id"]; ok {
		userID = int(u.(float64))
	}

	if l, ok := m["login"]; ok {
		login = l.(string)
	}

	return userID, login, nil
}

func (auth *Auth) GetNewToken(userID int, login string) (string, error) {
	_, tokenString, err := auth.tokenAuth.Encode(map[string]interface{}{"user_id": userID, "login": login})
	return tokenString, err
}

func (auth *Auth) GetJWTAuth() *jwtauth.JWTAuth {
	return auth.tokenAuth
}

func bakeCookie(token string) (*http.Cookie, error) {
	var c http.Cookie
	c.Name = "jwt"
	c.Value = token
	return &c, nil
}
