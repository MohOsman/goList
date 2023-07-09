package storage

import "goList/types"

type Memeory struct {
	
}

func  NewMemory() *Memeory{
	return  &Memeory{}
}

func(s *Memeory) Get(id int ) *types.User  {
		return &types.User{
		ID: 1,
		Username: "user1",
		Password: "password",
		}
}