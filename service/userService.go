package service

import (
	"errors"
	"goList/types"
	"goList/storage"
	"goList/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	// Check if the username doesn't contain special characters
	if !utils.IsUsernameValid(user.Username) {
		err := errors.New("Username is not a valid username")
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error while hashing password: %v", err)
		return err
	}

	hashedUser := types.UserDAO{
		Username: user.Username,
		Password: hashedPassword,
	}

	err = us.s.RegisterUser(hashedUser)
	if err != nil {
		log.Printf("UserService: Error creating user %v", err)
		return err
	}
	return nil
}

func (us *UserService) Login(user types.User) (*string, error) {
	// // fetch with username
	userDAO, err := us.s.FindUserByUsername(user.Username)
	if err != nil {
		log.Printf("UserService: Error getting user %v", err)
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(userDAO.Password, []byte(user.Password))
	if err != nil {
		log.Printf("Password not matching")
		return nil, err
	}

	return &user.Username, nil

}
