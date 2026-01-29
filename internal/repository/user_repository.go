package repository

import (
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/utils/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	db := database.GetInstance()
	return db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	db := database.GetInstance()
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	db := database.GetInstance()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) AddAddress(address *model.Address) error {
	db := database.GetInstance()
	return db.Create(address).Error
}

func (r *UserRepository) GetAddresses(userID uint) ([]model.Address, error) {
	db := database.GetInstance()
	var addresses []model.Address
	if err := db.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}
