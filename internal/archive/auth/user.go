package auth

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ARTM2000/archive1/internal/archive/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func NewUserManager(config UserConfig, userRepo database.UserRepository) userManger {
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
	userRepository database.UserRepository
	config         UserConfig
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

func (um *userManger) LoginUser(email string, password string) (string, error) {
	user, err := um.userRepository.FindUserWithEmail(email)
	if err != nil {
		log.Default().Println("error in finding user with email", err.Error())
		if errors.Is(err, database.ErrRecordNotFound) {
			return "", ErrEmailOrPasswordIsIncorrect
		}
		return "", ErrUnhandled
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		log.Default().Println("error in comparing password in login", err.Error())
		return "", ErrEmailOrPasswordIsIncorrect
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
		return "", ErrUnhandled
	}

	return tokenString, nil
}

func (um *userManger) VerifyUserAccessToken(token string) (*database.UserSchema, error) {
	tokenByte, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return []byte(um.config.JWTSecret), nil
	})

	if err != nil {
		log.Default().Println("error in parsing access token.", err.Error())
		return nil, ErrUnauthorized
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok {
		log.Default().Println("can not retrieve claims from token")
		return nil, ErrUnauthorized
	}

	ext := claims["ext"].(map[string]interface{})
	userIdStr := ext["id"].(string)
	userId, _ := strconv.ParseUint(userIdStr, 10, 64)

	user, err := um.userRepository.FindUserWithId(uint(userId))
	if err != nil {
		log.Default().Println("error in retrieving user from database", err.Error())
		return nil, ErrUnauthorized
	}

	return user, nil
}
