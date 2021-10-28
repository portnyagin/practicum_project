package handler

import (
	"encoding/json"
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
