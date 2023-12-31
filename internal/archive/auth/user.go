package auth

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ARTM2000/archivo/internal/archive/xerrors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func NewUserManager(config UserConfig, userRepo UserRepository) userManger {
	return userManger{
		userRepository: userRepo,
		config:         config,
	}
}

type UserConfig struct {
	JWTSecret     string
	JWTExpireTime time.Duration
}

type userManger struct {
	userRepository UserRepository
	config         UserConfig
}

func (um *userManger) AdminExistenceCheck() (bool, error) {
	adminUser, err := um.userRepository.FindAdminUser()
	if err != nil && !errors.Is(err, xerrors.ErrRecordNotFound) {
		log.Default().Printf("error in finding admin user, error: %s", err.Error())
		return false, err
	}

	if adminUser != nil {
		log.Default().Println("admin user exist")
		return true, nil
	}

	log.Default().Println("admin user not exist")
	return false, nil
}

func (um *userManger) RegisterAdmin(email string, username string, password string) (*User, error) {
	adminUser, err := um.userRepository.FindAdminUser()
	if adminUser != nil {
		log.Default().Println("admin user exists")
		return nil, xerrors.ErrAdminExist
	}

	if err != nil && !errors.Is(err, xerrors.ErrRecordNotFound) {
		log.Default().Println("unhandled error occurs.", err.Error())
		return nil, xerrors.ErrUnhandled
	}

	// in case that no admin user exists, create new one
	passwordByte := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		log.Default().Println("password hashing problem.", err.Error())
		return nil, xerrors.ErrUnhandled
	}

	newAdminUser, err := um.userRepository.CreateNewAdminUser(username, email, string(passwordHash))
	if err != nil {
		if errors.Is(err, xerrors.ErrDuplicateViolation) {
			return nil, xerrors.ErrEmailOrUsernameInUse
		}

		log.Default().Println("create admin user error.", err.Error())
		return nil, xerrors.ErrUnhandled
	}

	return newAdminUser, nil
}

func (um *userManger) RegisterUser(email string, username string, password string) (*User, error) {
	existingUser, err := um.userRepository.FindUserWithEmailOrUsername(email, username)
	if err != nil && !errors.Is(err, xerrors.ErrRecordNotFound) {
		log.Default().Printf("[Unhandled] error in check user existence with same username or password. error: %s", err.Error())
		return nil, xerrors.ErrUnhandled
	}
	if existingUser != nil {
		log.Default().Println("user with this email or username exists")
		return nil, xerrors.ErrUserExist
	}

	passwordByte := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		log.Default().Println("password hashing problem.", err.Error())
		return nil, xerrors.ErrUnhandled
	}

	newNonAdminUser, err := um.userRepository.CreateNewNonAdminUser(username, email, string(passwordHash))
	if err != nil {
		if errors.Is(err, xerrors.ErrDuplicateViolation) {
			return nil, xerrors.ErrEmailOrUsernameInUse
		}

		log.Default().Println("create non admin user error.", err.Error())
		return nil, xerrors.ErrUnhandled
	}

	return newNonAdminUser, nil
}

func (um *userManger) LoginUser(email string, password string) (string, error) {
	user, err := um.userRepository.FindUserWithEmail(email)
	if err != nil {
		log.Default().Println("error in finding user with email", err.Error())
		if errors.Is(err, xerrors.ErrRecordNotFound) {
			return "", xerrors.ErrEmailOrPasswordIsIncorrect
		}
		return "", xerrors.ErrUnhandled
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		log.Default().Println("error in comparing password in login", err.Error())
		return "", xerrors.ErrEmailOrPasswordIsIncorrect
	}

	now := time.Now().UTC()
	claims := &jwt.MapClaims{
		"exp": now.Add(um.config.JWTExpireTime).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"ext": map[string]string{
			"id": fmt.Sprint(user.ID),
		},
	}

	accessTokenByte := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessTokenByte.SignedString([]byte(um.config.JWTSecret))

	if err != nil {
		fmt.Println("error in creating token", err.Error())
		return "", xerrors.ErrUnhandled
	}

	um.userRepository.UpdateLastLoginTime(user.ID)
	if err != nil {
		fmt.Println("error in updating last login time", err.Error())
		return "", xerrors.ErrUnhandled
	}

	return tokenString, nil
}

func (um *userManger) ChangeInitialPassword(email, initialPassword, newPassword string) error {
	user, err := um.userRepository.FindUserWithEmail(email)
	if err != nil {
		log.Default().Println("error in finding user with email", err.Error())
		if errors.Is(err, xerrors.ErrRecordNotFound) {
			return xerrors.ErrEmailOrPasswordIsIncorrect
		}
		return xerrors.ErrUnhandled
	}

	if !user.ChangeInitialPassword {
		return xerrors.ErrUserInitialPasswordHasBeenChanged
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(initialPassword))
	if err != nil {
		log.Default().Println("error in comparing password in change initial password", err.Error())
		return xerrors.ErrUnauthorized
	}

	passwordByte := []byte(newPassword)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		log.Default().Println("password hashing problem.", err.Error())
		return xerrors.ErrUnauthorized
	}

	_, err = um.userRepository.ChangeUserPassword(user.ID, string(passwordHash))
	if err != nil {
		log.Default().Println("error in changing user password in change initial password", err)
		return xerrors.ErrUnhandled
	}

	return nil
}

func (um *userManger) VerifyUserAccessToken(token string) (*User, error) {
	tokenByte, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return []byte(um.config.JWTSecret), nil
	})

	if err != nil {
		log.Default().Println("error in parsing access token.", err.Error())
		return nil, xerrors.ErrUnauthorized
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok {
		log.Default().Println("can not retrieve claims from token")
		return nil, xerrors.ErrUnauthorized
	}

	ext := claims["ext"].(map[string]interface{})
	userIdStr := ext["id"].(string)
	userId, _ := strconv.ParseUint(userIdStr, 10, 64)

	user, err := um.userRepository.FindUserWithId(uint(userId))
	if err != nil {
		log.Default().Println("error in retrieving user from database", err.Error())
		return nil, xerrors.ErrUnauthorized
	}

	return user, nil
}

func (um *userManger) IsUserAdmin(user *User) bool {
	return user.IsAdmin
}

func (um *userManger) GetAllUsers(option FindAllOption) (*[]User, int64, error) {
	usersList, totalUsers, err := um.userRepository.FindAllUsers(option)

	if err != nil {
		log.Default().Println("[Unhandled] error in finding all source servers", err.Error())
		return nil, 0, xerrors.ErrUnhandled
	}

	return usersList, totalUsers, nil
}
