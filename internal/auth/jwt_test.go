package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateAndVerify(t *testing.T) {
	t.Parallel()

	manager := NewJWTManager("test-secret", time.Minute)

	token, err := manager.Generate("testuser", "user-123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.Verify(token)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "user-123", claims.UserID)
}

func TestJWTManager_Verify_InvalidCases(t *testing.T) {
	t.Parallel()

	manager := NewJWTManager("test-secret", time.Minute)

	// Generate valid token for tampering
	validToken, err := manager.Generate("testuser", "user-123")
	assert.NoError(t, err)

	// Tampered token
	tamperedToken := validToken + "tamper"

	// Expired token
	expiredManager := NewJWTManager("test-secret", -time.Minute) // already expired
	expiredToken, err := expiredManager.Generate("testuser", "user-123")
	assert.NoError(t, err)

	cases := []struct {
		name       string
		token      string
		expectErr  bool
	}{
		{"invalid token format", "invalid.token.here", true},
		{"tampered token", tamperedToken, true},
		{"expired token", expiredToken, true},
		{"valid token", validToken, false},
	}

	for _, tc := range cases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			claims, err := manager.Verify(tc.token)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "testuser", claims.Username)
				assert.Equal(t, "user-123", claims.UserID)
			}
		})
	}
}
