package archive

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"github.com/ARTM2000/archivo/internal/archive/sourceserver"
	"github.com/ARTM2000/archivo/internal/archive/xerrors"
	"github.com/ARTM2000/archivo/internal/validate"
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
	FileName string                `form:"filename" validate:"omitempty,filename,alphanum"`
	Rotate   int                   `form:"rotate" validate:"required,number"`
}

type listData struct {
	SortBy    string `query:"sort_by" validate:"required"`
	SortOrder string `query:"sort_order" validate:"required"`
	Start     *int   `query:"start" validate:"omitempty,number"`
	End       *int   `query:"end" validate:"omitempty,number"`
}

type snapshotListData struct {
	SrvId    uint   `params:"srvId" validate:"required,number"`
	Filename string `params:"filename" validate:"required"`
}

type downloadSnapshotData struct {
	snapshotListData
	Snapshot string `params:"snapshot" validate:"required"`
}

func (api *API) getListOfSourceServers(c *fiber.Ctx) error {
	var data listData
	if err := c.QueryParser(&data); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs, ok := validate.ValidateStruct[listData](&data); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	srcsrvManager := sourceserver.NewSrvManager(
		sourceserver.SrvConfig{},
		sourceserver.NewSrvRepository(api.DB),
	)

	if data.Start == nil {
		var initialStart = 0
		data.Start = &initialStart
	}
	if data.End == nil {
		var initialEnd = 10
		data.End = &initialEnd
	}

	servers, total, err := srcsrvManager.GetListOfAllSourceServers(sourceserver.FindAllOption{
		SortBy:    data.SortBy,
		SortOrder: data.SortOrder,
		Start:     *data.Start,
		End:       *data.End,
	})

	if err != nil {
		log.Default().Println("[Unhandled] error for finding all source servers", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: map[string]interface{}{
			"list":  servers,
			"total": total,
		},
	}))
}

func (api *API) getSourceServerFilesList(c *fiber.Ctx) error {
	params := struct {
		SrvId uint `params:"srvId" validate:"required,number"`
	}{}
	if err := c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errs, ok := validate.ValidateStruct[struct {
		SrvId uint `params:"srvId" validate:"required,number"`
	}](&params); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	var lData listData
	if err := c.QueryParser(&lData); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errs, ok := validate.ValidateStruct[listData](&lData); !ok {
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

	if lData.Start == nil {
		var initialStart = 0
		lData.Start = &initialStart
	}
	if lData.End == nil {
		var initialEnd = 10
		lData.End = &initialEnd
	}

	filesList, total, err := srcsrvManager.GetListOfSourceServerFiles(
		params.SrvId,
		sourceserver.FindAllOption{
			SortBy:    lData.SortBy,
			SortOrder: lData.SortOrder,
			Start:     *lData.Start,
			End:       *lData.End,
		},
	)
	if err != nil {
		log.Default().Printf("error in getting files list of source server by id '%d'. error: %s", params.SrvId, err.Error())
		if errors.Is(err, xerrors.ErrNoStoreForSourceServer) {
			return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
				Data: map[string]interface{}{
					"list":  filesList,
					"total": total,
				},
			}))
		}

		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: map[string]interface{}{
			"list":  filesList,
			"total": total,
		},
	}))
}

func (api *API) getListOfFileSnapshots(c *fiber.Ctx) error {
	params := snapshotListData{}
	if err := c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errs, ok := validate.ValidateStruct[snapshotListData](&params); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	var lData listData
	if err := c.QueryParser(&lData); err != nil {
		log.Default().Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errs, ok := validate.ValidateStruct[listData](&lData); !ok {
		log.Default().Println(errs[0].Message)
		return fiber.NewError(fiber.StatusUnprocessableEntity, errs[0].Message)
	}

	if lData.Start == nil {
		var initialStart = 0
		lData.Start = &initialStart
	}
	if lData.End == nil {
		var initialEnd = 10
		lData.End = &initialEnd
	}

	srcsrvManager := sourceserver.NewSrvManager(
		sourceserver.SrvConfig{
			CorrelationId:   c.GetRespHeader(fiber.HeaderXRequestID),
			StoreMode:       api.Config.FileStore.Mode,
			DiskStoreConfig: sourceserver.DiskStoreConfig(api.Config.FileStore.DiskConfig),
		},
		sourceserver.NewSrvRepository(api.DB),
	)

	snapshotsList, total, err := srcsrvManager.GetListOfFileSnapshotsByFilenameAndSrvId(
		params.SrvId,
		params.Filename,
		sourceserver.FindAllOption{
			SortBy:    lData.SortBy,
			SortOrder: lData.SortOrder,
			Start:     *lData.Start,
			End:       *lData.End,
		},
	)

	if err != nil {
		log.Default().Printf("error in getting files list of source server by id '%d' and filename '%s'. error: %s", params.SrvId, params.Filename, err.Error())
		if errors.Is(err, xerrors.ErrNoStoreForSourceServer) || errors.Is(err, xerrors.ErrNoFileStoredOnSourceServerByThisName) {
			return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
				Data: map[string]interface{}{
					"list":  snapshotsList,
					"total": total,
				},
			}))
		}

		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: map[string]interface{}{
			"list":  snapshotsList,
			"total": total,
		},
	}))
}

func (api *API) downloadSnapshot(c *fiber.Ctx) error {
	params := downloadSnapshotData{}
	if err := c.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errs, ok := validate.ValidateStruct[downloadSnapshotData](&params); !ok {
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

	snapshotByte, filename, err := srcsrvManager.ReadSnapshot(params.SrvId, params.Filename, params.Snapshot)
	if err != nil {
		if errors.Is(err, xerrors.ErrSnapshotNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	c.Append(fiber.HeaderContentType, "application/octet-stream")
	c.Append(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
	_, err = c.Status(fiber.StatusOK).Write(*snapshotByte)
	if err != nil {
		log.Default().Printf("error in writing file to response, error: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	return nil
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
			"id":      newSourceServerD.NewServer.ID,
			"name":    newSourceServerD.NewServer.Name,
			"api_key": newSourceServerD.APIKey,
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
		if errors.Is(err, xerrors.ErrFileRotateCountIsLowerThanPreviousOne) ||
			errors.Is(err, xerrors.ErrRotateGlobalLimitReached) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Message: "done",
	}))
}

func (api *API) storeCommonStatistics(c *fiber.Ctx) error {
	srcsrvManager := sourceserver.NewSrvManager(
		sourceserver.SrvConfig{
			CorrelationId:   c.GetRespHeader(fiber.HeaderXRequestID),
			StoreMode:       api.Config.FileStore.Mode,
			DiskStoreConfig: sourceserver.DiskStoreConfig(api.Config.FileStore.DiskConfig),
		},
		sourceserver.NewSrvRepository(api.DB),
	)

	sourceServersCount, err := srcsrvManager.SourceServersCount()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	filesForBackupCount, err := srcsrvManager.SourceServerFilesCount()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	totalSnapshotOccupiedSize, err := srcsrvManager.TotalSnapshotsSize()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
		Data: map[string]interface{}{
			"backup_files_count":     filesForBackupCount,
			"source_servers_count":   sourceServersCount,
			"snapshot_occupied_size": totalSnapshotOccupiedSize,
		},
	}))
}
