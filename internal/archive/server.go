package archive

import (
	"errors"
	"fmt"
	"log"

	"github.com/ARTM2000/archive1/web"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
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
	sessionStore := session.New(session.Config{
		Expiration: c.Auth.JWTExpireTime,
	})

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
		Config:       c,
		SessionStore: sessionStore,
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
		AllowOrigins:     "http://localhost:5173,http://127.0.0.1:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		ExposeHeaders:    "Content-Length,Content-Disposition,Content-Type",
		AllowCredentials: true,
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
			rt.Post("/logout", api.logoutUser)

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
			rtr.Get("/", api.getListOfSourceServers)
			rtr.Post("/new", api.registerNewSourceServer)
			rtr.Get("/:srvId/files", api.getSourceServerFilesList)
			rtr.Get("/:srvId/files/:filename", api.getListOfFileSnapshots)
			rtr.Get("/:srvId/files/:filename/:snapshot/download", api.downloadSnapshot)
		})

		router.Route("/users", func(rtr fiber.Router) {
			rtr.Use(api.authorizationMiddleware)
			// admin only
			rtr.Use(api.adminAuthorizationMiddleware)
			rtr.Get("/", api.getAllUsersInformation)
		})
	}, "APIv1")

	app.Use(web.ServePath, web.ServeDashboard)
	app.Use("/", func(c *fiber.Ctx) error {
		return c.Redirect(web.ServePath, fiber.StatusTemporaryRedirect)
	})

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
	DB           *gorm.DB
	Config       *Config
	SessionStore *session.Store
}
