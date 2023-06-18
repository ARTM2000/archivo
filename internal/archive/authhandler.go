package archive

import (
	"errors"
	"log"

	"github.com/ARTM2000/archive1/internal/archive/auth"
	"github.com/ARTM2000/archive1/internal/validate"
	"github.com/gofiber/fiber/v2"
)

type registerAdminDto struct {
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

	userManger := auth.NewUserManager(auth.UserConfig{}, api.DBM.NewUserRepository())
	newAdmin, err := userManger.RegisterAdmin(registerData.Email, registerData.Username, registerData.Password)
	if err != nil {
		if errors.Is(err, auth.ErrAdminExist) {
			log.Default().Println("admin exists")
			return fiber.NewError(fiber.StatusConflict, "admin already exists")
		}

		log.Default().Println("unhandled error from adminManager")
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: fiber.Map{
			"admin": newAdmin,
		},
		Message: "new admin user created",
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
		api.DBM.NewUserRepository(),
	)
	token, err := userManager.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		if errors.Is(err, auth.ErrEmailOrPasswordIsIncorrect) {
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
