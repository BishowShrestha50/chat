// services/auth_service.go
package service

import (
	"chat/model"
	"chat/utils"
	"errors"

	"gorm.io/gorm"
)

type IAuthService interface {
	Register(credentials model.Credentials) error
	Login(credentials model.Credentials) (string, error)
	GetUserById(id uint) (model.User, error)
}

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) IAuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(credential model.Credentials) error {
	hashedPassword, err := utils.HashPassword(credential.Password)
	if err != nil {
		return err
	}

	user := model.User{Username: credential.Username, Password: hashedPassword}
	return s.DB.Create(&user).Error
}

func (s *AuthService) Login(credential model.Credentials) (string, error) {
	var user model.User
	if result := s.DB.Where("username = ?", credential.Username).First(&user); result.Error != nil {
		return "", result.Error
	}

	if !utils.CheckPasswordHash(credential.Password, user.Password) {
		return "", errors.New("invalid credentials")

	}

	return utils.GenerateJWT(user.ID)
}

func (s *AuthService) GetUserById(id uint) (model.User, error) {
	var user model.User
	if result := s.DB.First(&user, id); result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}
