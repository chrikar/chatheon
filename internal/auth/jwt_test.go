package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateAndVerify(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Minute)

	// Generate token
	token, err := manager.Generate("testuser", "user-123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	claims, err := manager.Verify(token)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "user-123", claims.UserID)
}

func TestJWTManager_Verify_InvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Minute)

	// Provide an invalid token
	_, err := manager.Verify("invalid.token.here")
	assert.Error(t, err)
}

func TestJWTManager_Verify_TamperedToken(t *testing.T) {
	manager := NewJWTManager("test-secret", time.Minute)

	// Generate token
	token, err := manager.Generate("testuser", "user-123")
	assert.NoError(t, err)

	// Tamper token
	tamperedToken := token + "tamper"

	_, err = manager.Verify(tamperedToken)
	assert.Error(t, err)
}

func TestJWTManager_Verify_ExpiredToken(t *testing.T) {
	manager := NewJWTManager("test-secret", -time.Minute) // already expired

	// Generate expired token
	token, err := manager.Generate("testuser", "user-123")
	assert.NoError(t, err)

	_, err = manager.Verify(token)
	assert.Error(t, err)
}
