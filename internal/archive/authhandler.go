package archive

import (
	"errors"
	"log"
	"strings"

	"github.com/ARTM2000/archive1/internal/archive/auth"
	"github.com/ARTM2000/archive1/internal/archive/xerrors"
	"github.com/ARTM2000/archive1/internal/validate"
	"github.com/gofiber/fiber/v2"
)

const (
	UserLocalName = "user"
)

type registerAdminDto struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,password"`
}

type registerUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,password"`
}

type loginUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (api *API) registerAdmin(c *fiber.Ctx) error {
	registerData := registerAdminDto{}
	if err := c.BodyParser(&registerData); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs, ok := validate.ValidateStruct[registerAdminDto](&registerData); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	userManger := auth.NewUserManager(auth.UserConfig{}, auth.NewUserRepository(api.DB))
	newAdmin, err := userManger.RegisterAdmin(registerData.Email, registerData.Username, registerData.Password)
	if err != nil {
		if errors.Is(err, xerrors.ErrAdminExist) {
			log.Default().Println("admin exists")
			return fiber.NewError(fiber.StatusConflict, "admin already exists")
		}

		log.Default().Println("unhandled error from userManager")
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: fiber.Map{
			"admin": newAdmin,
		},
		Message: "new admin user created",
	}))
}

func (api *API) registerUser(c *fiber.Ctx) error {
	registerData := registerUserDto{}
	if err := c.BodyParser(&registerData); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs, ok := validate.ValidateStruct[registerUserDto](&registerData); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	userManger := auth.NewUserManager(auth.UserConfig{}, auth.NewUserRepository(api.DB))
	newUser, err := userManger.RegisterUser(registerData.Email, registerData.Username, registerData.Password)
	if err != nil {
		if errors.Is(err, xerrors.ErrUserExist) {
			log.Default().Println("user (non admin) exists")
			return fiber.NewError(fiber.StatusConflict, "user with same email or username already exists")
		}

		log.Default().Println("unhandled error from userManager", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Message: "new user registered",
		Data: map[string]interface{}{
			"user": newUser,
		},
	}))
}

func (api *API) loginUser(c *fiber.Ctx) error {
	loginData := loginUserDto{}
	if err := c.BodyParser(&loginData); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userManager := auth.NewUserManager(
		auth.UserConfig{
			JWTSecret:     api.Config.Auth.JWTSecret,
			JWTExpireTime: api.Config.Auth.JWTExpireTime,
		},
		auth.NewUserRepository(api.DB),
	)
	token, err := userManager.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		if errors.Is(err, xerrors.ErrEmailOrPasswordIsIncorrect) {
			log.Default().Println("email or password is incorrect")
			return fiber.NewError(fiber.StatusUnauthorized, "email or password in incorrect")
		}
		log.Default().Println("unhandled error", err.Error())
		return fiber.NewError(fiber.StatusUnauthorized, "email or password in incorrect")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Message: "welcome",
		Data: map[string]interface{}{
			"access_token": token,
		},
	}))
}

func (api *API) authorizationMiddleware(c *fiber.Ctx) error {
	// todo: enable role base access control (RBAC)
	authHeader := c.Get(fiber.HeaderAuthorization)
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

	tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenStr == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

	userManager := auth.NewUserManager(
		auth.UserConfig{
			JWTSecret:     api.Config.Auth.JWTSecret,
			JWTExpireTime: api.Config.Auth.JWTExpireTime,
		},
		auth.NewUserRepository(api.DB),
	)

	user, err := userManager.VerifyUserAccessToken(tokenStr)
	if err != nil {
		if errors.Is(err, xerrors.ErrUnauthorized) {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
		}
		log.Default().Println("[Unhandled] error in verifying token", err.Error())
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

	c.Locals(UserLocalName, user)
	log.Default().Printf("request authorized. user: %+v \n", user)
	return c.Next()
}

func (api *API) getUserInfo(c *fiber.Ctx) error {
	user := c.Locals(UserLocalName).(*auth.User)

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: map[string]interface{}{
			"user": user,
		},
	}))
}
