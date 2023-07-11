package service

import (
	"goList/types"
	"goList/storage"
)

type UserService struct {
	s storage.UserStorage
}

func NewUserService(storage storage.UserStorage) *UserService {
	return &UserService{
		s: storage,
	}
}

func (us *UserService) RegisterUser(user types.User) error {
	return us.s.RegisterUser(user)
}
