package storage

import "goList/types"

type Storage interface {
	get( id int) *types.User
}