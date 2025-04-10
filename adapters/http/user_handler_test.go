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

func TestUserHandler_RegisterUser(t *testing.T) {
	service := new(mockUserService)
	handler := NewUserHandler(service)

	tests := []struct {
		name         string
		payload      registerRequest
		mockSetup    func()
		expectedCode int
	}{
		{
			name:    "valid",
			payload: registerRequest{"newuser", "password"},
			mockSetup: func() {
				service.On("Register", "newuser", "password").Return(nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:    "existing user",
			payload: registerRequest{"existing", "password"},
			mockSetup: func() {
				service.On("Register", "existing", "password").Return(application.ErrUsernameTaken)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service.ExpectedCalls = nil // reset calls
			tc.mockSetup()

			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			handler.RegisterUser(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
			service.AssertExpectations(t)
		})
	}
}

func TestUserHandler_LoginUser(t *testing.T) {
	service := new(mockUserService)
	handler := NewUserHandler(service)

	tests := []struct {
		name         string
		payload      loginRequest
		mockSetup    func()
		expectedCode int
		expectToken  bool
	}{
		{
			name:    "valid",
			payload: loginRequest{"user", "pass"},
			mockSetup: func() {
				service.On("Login", "user", "pass").Return("mock-token", nil)
			},
			expectedCode: http.StatusOK,
			expectToken:  true,
		},
		{
			name:    "invalid credentials",
			payload: loginRequest{"user", "wrongpass"},
			mockSetup: func() {
				service.On("Login", "user", "wrongpass").Return("", application.ErrInvalidCredentials)
			},
			expectedCode: http.StatusUnauthorized,
			expectToken:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service.ExpectedCalls = nil // reset calls
			tc.mockSetup()

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

			service.AssertExpectations(t)
		})
	}
}
