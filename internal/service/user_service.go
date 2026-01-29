package service

import (
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
)

type UserService struct {
	Repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.Repo.FindByID(userID)
}

func (s *UserService) AddAddress(userID uint, street, city, state, zip, country string) error {
	address := &model.Address{
		UserID:  userID,
		Street:  street,
		City:    city,
		State:   state,
		ZipCode: zip,
		Country: country,
	}
	return s.Repo.AddAddress(address)
}

func (s *UserService) GetAddresses(userID uint) ([]model.Address, error) {
	return s.Repo.GetAddresses(userID)
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.Repo.FindAllUsers()
}
