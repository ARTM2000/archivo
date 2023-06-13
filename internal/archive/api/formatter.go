package api

import "github.com/gofiber/fiber/v2"

type Data struct {
	Data    map[string]interface{}
	Message string
	IsError bool
}

type final struct {
	TrackId string                 `json:"track_id"`
	Error   bool                   `json:"error"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func FormatResponse(c *fiber.Ctx, resData Data) final {
	message := ""
	if resData.Message != "" {
		message = resData.Message
	}

	data := map[string]interface{}{}
	if resData.Data != nil {
		data = resData.Data
	}

	trackId := c.GetRespHeader(fiber.HeaderXRequestID)

	if trackId == "" {
		panic("track id should be defined in format response function")
	}

	return final{
		TrackId: trackId,
		Error:   resData.IsError,
		Message: message,
		Data:    data,
	}
}
