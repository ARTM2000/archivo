package archive

import (
	"errors"
	"fmt"
	"log"

	"github.com/ARTM2000/archive1/internal/archive/api"
	"github.com/ARTM2000/archive1/internal/archive/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func runServer(c *Config) {
	fConfig := fiber.Config{
		CaseSensitive:                true,
		ServerHeader:                 "none",
		AppName:                      "Archive1",
		DisablePreParseMultipartForm: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// check that if error was an fiber NewError and got status code,
			// specify that in error handler
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

			return c.Status(code).JSON(api.FormatResponse(c, api.Data{
				Message: err.Error(),
				IsError: true,
			}))
		},
	}
	app := fiber.New(fConfig)

	/**
	 * General configuration
	 */
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${pid}] '${ip}:${port}' ${status} - ${method} ${path}\n",
	}))
	app.Use(requestid.New())
	app.Use(helmet.New())

	app.Use(func(c *fiber.Ctx) error {
		contentType := c.Get("Content-Type")
		if c.Method() != "GET" && contentType != "application/json" && contentType != "multipart/form-data" {
			return fiber.NewError(fiber.StatusBadRequest, "Request body must be in 'application/json' or 'multipart/form-data' format")
		}
		return c.Next()
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(api.FormatResponse(c, api.Data{
			Message: "everything is fine",
		}))
	})

	app.Route("/api/v1", func(router fiber.Router) {
		router.Route("/auth", auth.Routes)
	}, "APIv1")

	port := 8010
	host := ""
	if c.ServerPort != nil {
		port = *c.ServerPort
	}
	if c.ServerHost != nil {
		host = *c.ServerHost
	}

	err := app.Listen(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalln(err.Error())
	}
}
