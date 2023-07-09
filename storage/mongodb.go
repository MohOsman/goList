package storage

import "goList/types"

type Mongo struct {}

// uppercase func name exports to other packages
func (s *Mongo) Get(id int ) * types.User{
	return &types.User{
		ID: 1,
		Username: "user1",
		Password: "password",
	}
}