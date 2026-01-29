package service_test

import (
	"errors"
	"testing"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository/mocks"
	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/beingaloksharma/book-backend/utils/crypto"
	"github.com/beingaloksharma/book-backend/utils/token"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestSignup(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	// Case 1: Success
	mockRepo.On("FindByEmail", "john@example.com").Return(nil, nil).Once()
	mockRepo.On("CreateUser", mock.Anything).Return(nil).Once()

	err := authService.Signup("John", "john@example.com", "password123", model.RoleUser)
	assert.NoError(t, err)

	// Case 2: User exists
	existingUser := &model.User{Email: "john@example.com"}
	mockRepo.On("FindByEmail", "john@example.com").Return(existingUser, nil).Once()

	err = authService.Signup("John", "john@example.com", "password123", model.RoleUser)
	assert.Error(t, err)
	assert.Equal(t, "user already exists", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	// Setup Token
	viper.Set("jwt.secret", "testsecret")
	token.Init()

	mockRepo := new(mocks.MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	password := "password123"
	hashedPwd, _ := crypto.HashPassword(password)
	user := &model.User{
		Model:    gorm.Model{ID: 1},
		Email:    "john@example.com",
		Password: hashedPwd,
		Role:     model.RoleUser,
	}

	// Case 1: Success
	mockRepo.On("FindByEmail", "john@example.com").Return(user, nil).Once()

	val, err := authService.Login("john@example.com", password)
	assert.NoError(t, err)
	assert.NotEmpty(t, val)

	// Case 2: User not found
	mockRepo.On("FindByEmail", "unknown@example.com").Return(nil, errors.New("not found")).Once()

	_, err = authService.Login("unknown@example.com", password)
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())

	// Case 3: Wrong password
	mockRepo.On("FindByEmail", "john@example.com").Return(user, nil).Once()

	_, err = authService.Login("john@example.com", "wrongpass")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())

	mockRepo.AssertExpectations(t)
}
