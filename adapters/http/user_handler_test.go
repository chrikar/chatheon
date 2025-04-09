package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/application"
)

type mockUserService struct{}

func (m *mockUserService) Register(username, password string) error {
	if username == "existing" {
		return application.ErrUsernameTaken
	}
	if username == "" {
		return application.ErrUsernameRequired
	}
	if password == "" {
		return application.ErrPasswordRequired
	}
	return nil
}

func (m *mockUserService) Login(username, password string) (string, error) {
	if username == "user" && password == "pass" {
		return "mock-token", nil
	}
	return "", application.ErrInvalidCredentials
}

func TestUserHandler_RegisterUser(t *testing.T) {
	handler := NewUserHandler(&mockUserService{})

	tests := []struct {
		name         string
		payload      registerRequest
		expectedCode int
	}{
		{"valid", registerRequest{"newuser", "password"}, http.StatusCreated},
		{"existing user", registerRequest{"existing", "password"}, http.StatusBadRequest},
		{"empty username", registerRequest{"", "password"}, http.StatusBadRequest},
		{"empty password", registerRequest{"user", ""}, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			handler.RegisterUser(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
		})
	}
}

func TestUserHandler_LoginUser(t *testing.T) {
	handler := NewUserHandler(&mockUserService{})

	tests := []struct {
		name         string
		payload      loginRequest
		expectedCode int
		expectToken  bool
	}{
		{"valid", loginRequest{"user", "pass"}, http.StatusOK, true},
		{"invalid credentials", loginRequest{"user", "wrongpass"}, http.StatusUnauthorized, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			handler.LoginUser(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
			if tc.expectToken {
				var resp loginResponse
				err := json.NewDecoder(rr.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Token)
			}
		})
	}
}
