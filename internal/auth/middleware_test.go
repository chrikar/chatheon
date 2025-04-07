package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTMiddleware(t *testing.T) {
	jwtManager := NewJWTManager("test-secret", time.Minute)

	// Generate valid token
	validToken, err := jwtManager.Generate("testuser", "user-123")
	assert.NoError(t, err)

	// Table-driven tests
	cases := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectUserID   string
		expectUsername string
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectUserID:   "user-123",
			expectUsername: "testuser",
		},
		{
			name:           "missing Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token format",
			authHeader:     "Bearer invalid.token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Dummy next handler to capture context values
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value(ContextUserIDKey)
				username := r.Context().Value(ContextUsernameKey)

				assert.Equal(t, tc.expectUserID, userID)
				assert.Equal(t, tc.expectUsername, username)

				w.WriteHeader(http.StatusOK)
			})

			// Build request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			// Record response
			rr := httptest.NewRecorder()

			// Apply middleware
			middleware := JWTMiddleware(jwtManager)
			middleware(nextHandler).ServeHTTP(rr, req)

			// Assert response status
			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
