package service_test

import (
	"errors"
	"testing"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository/mocks"
	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProfile(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	user := &model.User{Name: "John"}
	mockRepo.On("FindByID", uint(1)).Return(user, nil)

	result, err := userService.GetProfile(1)
	assert.NoError(t, err)
	assert.Equal(t, "John", result.Name)

	mockRepo.AssertExpectations(t)
}

func TestGetProfile_Error(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("db error"))

	_, err := userService.GetProfile(1)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestAddAddress(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	mockRepo.On("AddAddress", mock.AnythingOfType("*model.Address")).Return(nil)

	err := userService.AddAddress(1, "Street", "City", "State", "Zip", "Country")
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAddAddress_Error(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	mockRepo.On("AddAddress", mock.AnythingOfType("*model.Address")).Return(errors.New("db error"))

	err := userService.AddAddress(1, "Street", "City", "State", "Zip", "Country")
	assert.Error(t, err)
}

func TestGetAddresses(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	addresses := []model.Address{{City: "City"}}
	mockRepo.On("GetAddresses", uint(1)).Return(addresses, nil)

	result, err := userService.GetAddresses(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))

	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockRepo)

	users := []model.User{{Name: "John"}}
	mockRepo.On("FindAllUsers").Return(users, nil)

	result, err := userService.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))

	mockRepo.AssertExpectations(t)
}
