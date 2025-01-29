package services

import (
	"time"

	"github.com/bytetiff/bytechat/internal/models"
	repositories "github.com/bytetiff/bytechat/internal/repos"
	"github.com/bytetiff/bytechat/internal/utils"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser - регистрация нового пользователя (хешируем пароль)
func (s *UserService) RegisterUser(username, password string) (models.User, error) {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Username:  username,
		Password:  hashed,
		CreatedAt: time.Now(),
	}
	return s.repo.CreateUser(user)
}

// GetUserByUsername - получить пользователя по логину (для логина)
func (s *UserService) GetUserByUsername(username string) (models.User, error) {
	return s.repo.GetUserByUsername(username)
}
