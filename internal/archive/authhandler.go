package archive

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ARTM2000/archivo/internal/archive/auth"
	"github.com/ARTM2000/archivo/internal/archive/xerrors"
	"github.com/ARTM2000/archivo/internal/validate"
	"github.com/gofiber/fiber/v2"
)

const (
	UserLocalName        = "user"
	SessionCredentialKey = "tkn"
)

type registerAdminDto struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,password"`
}

type registerUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,alphanum,min=8"` // as this password acts as initial password, we will keep it simple
}

type loginUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type changeInitialPassword struct {
	InitialPassword string `json:"initial_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,password"`
}

type userActivityParams struct {
	UserID uint `params:"userId" validate:"required,numeric"`
}

func (api *API) checkAdminExistence(c *fiber.Ctx) error {
	userManger := auth.NewUserManager(auth.UserConfig{}, auth.NewUserRepository(api.DB))
	adminExist, err := userManger.AdminExistenceCheck()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: map[string]interface{}{
			"admin_exist": adminExist,
		},
	}))
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

		if errors.Is(err, xerrors.ErrEmailOrUsernameInUse) {
			log.Default().Println("admin with this email for username exist")
			return fiber.NewError(fiber.StatusConflict, "admin with this email for username exists")
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

func (api *API) changeUserInitialPassword(c *fiber.Ctx) error {
	changeInitPass := changeInitialPassword{}
	if err := c.BodyParser(&changeInitPass); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs, ok := validate.ValidateStruct[changeInitialPassword](&changeInitPass); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	userManger := auth.NewUserManager(auth.UserConfig{}, auth.NewUserRepository(api.DB))

	user := c.Locals(UserLocalName).(*auth.User)

	err := userManger.ChangeInitialPassword(user.Email, changeInitPass.InitialPassword, changeInitPass.NewPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, xerrors.ErrUnauthorized.Error())
	}

	session, err := api.SessionStore.Get(c)
	if err != nil {
		log.Default().Printf("error in getting session from store, error: %+v \n", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	if session.Get(SessionCredentialKey) != nil {
		log.Default().Println("delete session")
		session.Destroy()
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Message: "initial password changed",
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
			"id":   newUser.ID,
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

	session, err := api.SessionStore.Get(c)
	if err != nil {
		log.Default().Printf("error in getting session from store, error: %+v \n", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	session.Set(SessionCredentialKey, fmt.Sprintf("Bearer %s", token))
	err = session.Save()
	if err != nil {
		log.Default().Printf("error in saving session from store, error: %+v \n", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	log.Default().Println(session)

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Message: "welcome",
	}))
}

func (api *API) logoutUser(c *fiber.Ctx) error {
	session, err := api.SessionStore.Get(c)
	if err != nil {
		log.Default().Printf("error in getting session from store, error: %+v \n", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	if session.Get(SessionCredentialKey) != nil {
		log.Default().Println("delete session")
		session.Destroy()
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Message: "user logout was successful",
	}))
}

func (api *API) _commonAuthorization(c *fiber.Ctx) (*auth.User, error) {
	session, err := api.SessionStore.Get(c)
	if err != nil {
		log.Default().Printf("error in getting session from store, error: %+v \n", err.Error())
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	tknData := session.Get(SessionCredentialKey)
	if tknData == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}
	authHeader := tknData.(string)

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

	tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenStr == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
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
			return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
		}
		log.Default().Println("[Unhandled] error in verifying token", err.Error())
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

	return user, nil
}

func (api *API) preDashboardAuthorizationMiddleware(c *fiber.Ctx) error {
	user, err := api._commonAuthorization(c)
	if err != nil {
		return err
	}

	c.Locals(UserLocalName, user)
	log.Default().Printf("request authorized for pre dashboard actions. user: %+v \n", user)
	return c.Next()
}

func (api *API) authorizationMiddleware(c *fiber.Ctx) error {
	user, err := api._commonAuthorization(c)
	if err != nil {
		return err
	}

	if user.ChangeInitialPassword {
		// if user initial password did not changed, user is not authorized to
		// use dashboard and will force to change him/her password
		log.Default().Println("user should change him/her initial password to be authorized!")
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

	userActivityManager := auth.NewUserActivityManager(
		auth.NewUserActivityRepository(
			api.DB,
		),
	)
	err = userActivityManager.SaveNewActivity(user.ID, string(c.Request().Header.Method()), string(c.Request().RequestURI()))
	if err != nil {
		log.Default().Println("got error in user activity log >", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Locals(UserLocalName, user)
	log.Default().Printf("request authorized. user: %+v \n", user)
	return c.Next()
}

func (api *API) adminAuthorizationMiddleware(c *fiber.Ctx) error {
	user := c.Locals(UserLocalName).(*auth.User)

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}
	userManager := auth.NewUserManager(
		auth.UserConfig{
			JWTSecret:     api.Config.Auth.JWTSecret,
			JWTExpireTime: api.Config.Auth.JWTExpireTime,
		},
		auth.NewUserRepository(api.DB),
	)

	if isAdminPermitted := userManager.IsUserAdmin(user); !isAdminPermitted {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

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

func (api *API) getAllUsersInformation(c *fiber.Ctx) error {
	var data listData
	if err := c.QueryParser(&data); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs, ok := validate.ValidateStruct[listData](&data); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	userManager := auth.NewUserManager(
		auth.UserConfig{
			JWTSecret:     api.Config.Auth.JWTSecret,
			JWTExpireTime: api.Config.Auth.JWTExpireTime,
		},
		auth.NewUserRepository(api.DB),
	)

	if data.Start == nil {
		var initialStart = 0
		data.Start = &initialStart
	}
	if data.End == nil {
		var initialEnd = 10
		data.End = &initialEnd
	}

	users, totalUsers, err := userManager.GetAllUsers(auth.FindAllOption{
		SortBy:    data.SortBy,
		SortOrder: data.SortOrder,
		Start:     *data.Start,
		End:       *data.End,
	})

	if err != nil {
		log.Default().Println("[Unhandled] error for finding all users", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: map[string]interface{}{
			"list":  users,
			"total": totalUsers,
		},
	}))
}

func (api *API) getSingleUserActivities(c *fiber.Ctx) error {
	var data listData
	if err := c.QueryParser(&data); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs, ok := validate.ValidateStruct[listData](&data); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	params := userActivityParams{}
	if err := c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errs, ok := validate.ValidateStruct[userActivityParams](&params); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	if data.Start == nil {
		var initialStart = 0
		data.Start = &initialStart
	}
	if data.End == nil {
		var initialEnd = 10
		data.End = &initialEnd
	}

	userActivityManager := auth.NewUserActivityManager(auth.NewUserActivityRepository(api.DB))
	userActs, total, err := userActivityManager.GetListForSingleUser(params.UserID, auth.FindAllOption{
		SortBy:    data.SortBy,
		SortOrder: data.SortOrder,
		Start:     *data.Start,
		End:       *data.End,
	})

	if err != nil {
		log.Default().Println("[Unhandled] error for finding all user activities", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: map[string]interface{}{
			"list":  userActs,
			"total": total,
		},
	}))
}
