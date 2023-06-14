package auth

import (
	"errors"
	"log"

	"github.com/ARTM2000/archive1/internal/archive/database"
	"golang.org/x/crypto/bcrypt"
)

func NewUserManager(userRepo database.UserRepository) userManger {
	return userManger{
		userRepository: userRepo,
	}
}

type userManger struct {
	userRepository database.UserRepository
}

func (um *userManger) RegisterAdmin(email string, username string, password string) (*database.UserSchema, error) {
	adminUser, err := um.userRepository.FindAdminUser()
	if adminUser != nil {
		log.Default().Println("admin user exists")
		return nil, ErrAdminExist
	}

	if err != nil && !errors.Is(err, database.ErrRecordNotFound) {
		log.Default().Println("unhandled error occurs.", err.Error())
		return nil, ErrUnhandled
	}

	// in case that no admin user exists, create new one
	passwordByte := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		log.Default().Println("password hashing problem.", err.Error())
		return nil, ErrUnhandled
	}

	newAdminUser, err := um.userRepository.CreateNewAdminUser(username, email, string(passwordHash))
	if err != nil {
		if errors.Is(err, database.ErrDuplicateViolation) {
			return nil, ErrEmailOrUsernameInUse
		}

		log.Default().Println("create admin user error.", err.Error())
		return nil, ErrUnhandled
	}

	return newAdminUser, nil
}
