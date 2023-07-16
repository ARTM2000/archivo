package archive

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
)

func runServer(c *Config) {
	sConfig := fiber.Config{
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

			return c.Status(code).JSON(FormatResponse(c, Data{
				Message: err.Error(),
				IsError: true,
			}))
		},
	}
	app := fiber.New(sConfig)

	api := API{
		DB: NewDBConnection(DBConfig{
			DBHost:    c.Database.Host,
			DBPort:    c.Database.Port,
			DBUser:    c.Database.Username,
			DBPass:    c.Database.Password,
			DBName:    c.Database.Name,
			DBZone:    c.Database.Zone,
			DBSSLMode: c.Database.SSLMode,
		}),
		Config: c,
	}

	/**
	 * General configuration
	 */
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${pid}] '${ip}:${port}' ${status} - ${method} ${path}\n",
	}))
	app.Use(requestid.New(requestid.Config{
		Next: func(c *fiber.Ctx) bool {
			trackId := c.Get(fiber.HeaderXRequestID)
			if trackId != "" {
				c.Set(fiber.HeaderXRequestID, trackId)
				return true
			}
			return false
		},
	}))
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Request-ID",
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Accepts(fiber.MIMEApplicationJSON, fiber.MIMEMultipartForm)
		return c.Next()
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(FormatResponse(c, Data{
			Message: "everything is fine",
		}))
	})

	app.Route("/api/v1", func(router fiber.Router) {
		router.Route("/auth", func(rt fiber.Router) {
			// non protected route
			rt.Get("/admin/existence", api.checkAdminExistence)
			rt.Post("/admin/register", api.registerAdmin)
			rt.Post("/user/register", api.registerUser)
			rt.Post("/login", api.loginUser)

			// protected routes
			rt.Use(api.authorizationMiddleware)
			rt.Get("/me", api.getUserInfo)
		})

		// protected routes
		router.Route("/servers", func(rtr fiber.Router) {
			rtr.Route("/store", func(rt fiber.Router) {
				rt.Use(api.authorizeSourceServerMiddleware)
				rt.Post("/file", api.rotateSrcSrvFile)
			})
			rtr.Use(api.authorizationMiddleware)
			rtr.Post("/new", api.registerNewSourceServer)
			rtr.Get("/list", api.getListOfSourceServers)
			rtr.Get("/:srvId/files", api.getSourceServerFilesList)
		})
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

// API handlers (controllers) register on this struct (class)
type API struct {
	DB     *gorm.DB
	Config *Config
}
