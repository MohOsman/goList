package storage

import "goList/types"

type UserStorage interface {
	RegisterUser(user types.UserDAO) error
	FindUserByUsername(username string) (*types.UserDAO, error)
}
