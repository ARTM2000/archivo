package archive

import (
	"errors"
	"log"

	"github.com/ARTM2000/archive1/internal/archive/sourceserver"
	"github.com/ARTM2000/archive1/internal/archive/xerrors"
	"github.com/ARTM2000/archive1/internal/validate"
	"github.com/gofiber/fiber/v2"
)

type registerNewSourceServer struct {
	Name string `json:"name" validate:"required,alphanum"`
}

func (api *API) registerNewSourceServer(c *fiber.Ctx) error {
	var registerData registerNewSourceServer
	if err := c.BodyParser(&registerData); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs, ok := validate.ValidateStruct[registerNewSourceServer](&registerData); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	srcsrvManager := sourceserver.NewSrvManager(sourceserver.NewSrvRepository(api.DB))
	newSourceServerD, err := srcsrvManager.RegisterNewSourceServer(registerData.Name)
	if err != nil {
		if errors.Is(err, xerrors.ErrSourceServerWithThisNameExists) {
			return fiber.NewError(fiber.StatusConflict, "source server with this name exists")
		}

		log.Default().Println("[Unhandled] error for registering new source server", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Message: "new source server created",
		Data: map[string]interface{}{
			"server": map[string]interface{}{
				"name":    newSourceServerD.NewServer.Name,
				"api_key": newSourceServerD.APIKey,
			},
		},
	}))
}
