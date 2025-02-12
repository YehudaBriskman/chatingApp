package services

import (
	"chatingApp/models"
	"chatingApp/repository"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.UserRepo.GetAllUsers()
}
