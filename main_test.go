package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	authService := NewAuthService()

	tests := []struct {
		name    string
		user    User
		wantErr bool
		wantMsg string
	}{
		{
			"valid user",
			User{"validuser", "user@example.com", "F9DXIK6hvuFINjmC"},
			false,
			"",
		},
		{
			"short password",
			User{"shortuser", "user@example.com", "F9DXIK"},
			true,
			"password must be at least 8 characters long",
		},
		{
			"invalid email",
			User{"invaliduser", "user@example", "F9DXIK6hvuFINjmC"},
			true,
			"invalid email format",
		},
		{
			"nickname already exists",
			User{"user", "user@example.com", "F9DXIK6hvuFINjmC"},
			true,
			"nickname already exists",
		},
		{
			"email already exists",
			User{"user", "test@example.com", "F9DXIK6hvuFINjmC"},
			true,
			"email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Authenticate(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("wanted error: %v, got: %v", tt.wantErr, err != nil)
			}

			if err != nil && err.Error() != tt.wantMsg {
				t.Errorf("wanted error message: %v, got: %v", tt.wantMsg, err.Error())
			}
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	authService := NewAuthService()
	handler := registerHandler(authService)

	tests := []struct {
		name     string
		method   string
		user     User
		wantCode int
		wantBody string
	}{
		{
			"successful registration",
			http.MethodPost,
			User{"newUser", "newuser@example.com", "F9DXIK6hvuFINjmC"},
			http.StatusCreated,
			"User registered successfully\n",
		},
		{
			"method not allowed",
			http.MethodGet,
			User{},
			http.StatusMethodNotAllowed,
			"method not allowed\n",
		},
		{
			"duplicate nickname",
			http.MethodPost,
			User{"user", "user@example.com", "F9DXIK6hvuFINjmC"},
			http.StatusBadRequest,
			"nickname already exists\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.user)
			r := httptest.NewRequest(tt.method, "/register", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("wanted status code: %v, got: %v", tt.wantCode, w.Code)
			}
			if w.Body.String() != tt.wantBody {
				t.Errorf("wanted body: %q, got: %q", tt.wantBody, w.Body.String())
			}
		})
	}
}
