package archive

import (
	"errors"
	"log"
	"mime/multipart"
	"strings"

	"github.com/ARTM2000/archive1/internal/archive/sourceserver"
	"github.com/ARTM2000/archive1/internal/archive/xerrors"
	"github.com/ARTM2000/archive1/internal/validate"
	"github.com/gofiber/fiber/v2"
)

const (
	SrcSrvLocalName = "srcsrv"
)

type registerNewSourceServer struct {
	Name string `json:"name" validate:"required,alphanum"`
}

type rotateSrcSrvFile struct {
	File     *multipart.FileHeader `form:"file" validate:"required"`
	FileName string                `form:"filename" validate:"omitempty,alphanum"`
	Rotate   int                   `form:"rotate" validate:"required,number"`
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

	srcsrvManager := sourceserver.NewSrvManager(
		sourceserver.SrvConfig{},
		sourceserver.NewSrvRepository(api.DB),
	)
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

func (api *API) authorizeSourceServerMiddleware(c *fiber.Ctx) error {
	sourceServerName := c.Get("X-Agent1-Name")
	if strings.TrimSpace(sourceServerName) == "" {
		log.Default().Println("unauthorized agent1 request. no agent name received")
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

	authHeader := c.Get(fiber.HeaderAuthorization)
	if strings.TrimSpace(authHeader) == "" {
		log.Default().Println("unauthorized agent1 request. no api key received")
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

	srcsrvManager := sourceserver.NewSrvManager(
		sourceserver.SrvConfig{},
		sourceserver.NewSrvRepository(api.DB),
	)
	srcSrv, err := srcsrvManager.AuthorizeSourceServer(sourceServerName, authHeader)
	if err != nil {
		log.Default().Printf("error in authorizing agent request, %s", err.Error())
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized request")
	}

	c.Locals(SrcSrvLocalName, srcSrv)
	return c.Next()
}

func (api *API) rotateSrcSrvFile(c *fiber.Ctx) error {
	var rotateData rotateSrcSrvFile
	if err := c.BodyParser(&rotateData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var fileLoadErr error
	rotateData.File, fileLoadErr = c.FormFile("file")
	if fileLoadErr != nil {
		log.Default().Println(fileLoadErr.Error())
	}

	if errs, ok := validate.ValidateStruct[rotateSrcSrvFile](&rotateData); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	srcsrvManager := sourceserver.NewSrvManager(
		sourceserver.SrvConfig{
			CorrelationId:   c.GetRespHeader(fiber.HeaderXRequestID),
			StoreMode:       api.Config.FileStore.Mode,
			DiskStoreConfig: sourceserver.DiskStoreConfig(api.Config.FileStore.DiskConfig),
		},
		sourceserver.NewSrvRepository(api.DB),
	)

	srcsrv := c.Locals(SrcSrvLocalName).(*sourceserver.SourceServer)

	err := srcsrvManager.RotateFile(srcsrv, rotateData.Rotate, rotateData.FileName, rotateData.File)
	if err != nil {
		log.Default().Println("error in file rotation. error:", err.Error())
		if errors.Is(err, xerrors.ErrFileRotateCountIsLowerThanPreviousOne) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Message: "done",
	}))
}