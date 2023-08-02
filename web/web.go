package web

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

//go:embed dist/*
var dashboardUI embed.FS

const ServePath = "/panel/"

func getFS() fs.FS {
	dist, err := fs.Sub(dashboardUI, "dist")
	if err != nil {
		// This can't happen... Go would throw a compilation error.
		panic(err)
	}
	return dist
}

func ServeDashboard(c *fiber.Ctx) error {
	if c.Method() != "GET" {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "method not allowed")
	}

	uiFS := getFS()
	path := c.Path()

	if path == "/panel/" {
		path = "index.html"
	}
	path = strings.TrimPrefix(path, "/panel/")

	file, err := uiFS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Default().Println("file", path, "not found:", err)
			return fiber.NewError(fiber.StatusNotFound, "page not found")
		}

		log.Default().Println("file", path, "cannot be read:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	c.Set(fiber.HeaderContentType, contentType)
	if strings.HasPrefix(path, "assets/") {
		c.Set("Cache-Control", "public, max-age=31536000")
	}

	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		c.Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	bs := make([]byte, stat.Size())
	bufio.NewReader(file).Read(bs)

	_, err = c.Status(fiber.StatusOK).Write(bs)
	return err
}
