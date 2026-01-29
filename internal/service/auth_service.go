package service

import (
	"errors"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
	"github.com/beingaloksharma/book-backend/utils/crypto"
	"github.com/beingaloksharma/book-backend/utils/token"
)

type AuthService struct {
	Repo repository.UserRepositoryInterface
}

func NewAuthService(repo repository.UserRepositoryInterface) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Signup(name, email, password string, role model.Role) error {
	existing, _ := s.Repo.FindByEmail(email)
	if existing != nil {
		return errors.New("user already exists")
	}

	hashedPwd, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}

	// Validate Role
	if role == "" {
		role = model.RoleUser
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: hashedPwd,
		Role:     role,
	}

	return s.Repo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !crypto.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return token.GenerateToken(user.ID, string(user.Role))
}
