package service

import (
	db "appContract/pkg/db/repository"
)

type UserService struct{}

func (s *UserService) Authorize(login string, password string) (int, error) {
	return db.Authorize(login, password)
}

func (s *UserService) GetAdmin(id int) (bool, error) {
	return db.GetAddmin(id)
}

func (s *UserService) ChangePassword(id int, password string) error {
	return db.ChangePassword(id, password)
}
