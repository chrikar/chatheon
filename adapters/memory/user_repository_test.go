package memory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/domain"
)

func TestUserRepository_CreateAndFind(t *testing.T) {
	t.Parallel()

	repo := NewUserRepository()

	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}

	// Test create
	err := repo.Create(user)
	assert.NoError(t, err)

	// Test duplicate
	err = repo.Create(user)
	assert.Error(t, err)

	// Test find existing
	foundUser, err := repo.FindByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user, foundUser)

	// Test find non-existing
	_, err = repo.FindByUsername("unknown")
	assert.Error(t, err)
}
