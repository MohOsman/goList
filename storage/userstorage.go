package storage

import "goList/types"

type UserStorage interface {
	RegisterUser(user types.User) error
}
