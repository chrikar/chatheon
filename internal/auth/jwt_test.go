package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateAndVerify(t *testing.T) {
	// 1h expiry for easy testing
	mgr := NewJWTManager("test-secret", time.Second)

	// Generate a token for user “alice”
	token, err := mgr.Generate("alice", "user-123")
	assert.NoError(t, err, "Generate should not error")

	// Immediately verify: should be valid
	claims, err := mgr.Verify(token)
	assert.NoError(t, err, "Verify should accept fresh token")
	assert.Equal(t, "alice", claims.Username, "Claims.Username should match")
	assert.Equal(t, "user-123", claims.UserID, "Claims.UserID should match")
	assert.NotEmpty(t, claims.ExpiresAt, "Claims.ExpiresAt should not be empty")

	// Tamper: invalid token should error
	_, err = mgr.Verify(token + "garbage")
	assert.Error(t, err, "Verify should reject malformed token")

	// Expiry: wait past 1s TTL
	time.Sleep(1100 * time.Millisecond)
	_, err = mgr.Verify(token)
	assert.Error(t, err, "Verify should reject expired token")
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
		name      string
		token     string
		expectErr bool
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
