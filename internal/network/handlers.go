package network

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mohammadne/mosaicman/pkg/utils"
)

var imageExtensions = []string{"jpg", "jpeg", "png"}

func (server *server) upload(c echo.Context) error {
	// Read file
	file, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "image file is missed")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	parts := strings.Split(file.Filename, ".")
	ext := parts[len(parts)-1]

	matched := false
	for _, imageExtension := range imageExtensions {
		if imageExtension == ext {
			matched = true
			break
		}
	}

	if !matched {
		return c.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
			Help    string `json:"help"`
		}{
			Message: "invalid image format has given",
			Help:    fmt.Sprintf("supported formats are: %v", imageExtensions),
		})
	}

	requestIP := c.Request().RemoteAddr
	uuid := utils.NewUUID(requestIP)

	fileName := fmt.Sprintf("%s.%s", uuid, ext)
	filePath := fmt.Sprintf("%s/%s", server.savePath, fileName)

	// Destination
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
		UUID    string `json:"uuid"`
	}{
		Message: "Image uploaded successfully",
		UUID:    uuid,
	})
}

func (server *server) generate(c echo.Context) error {
	return nil
}
