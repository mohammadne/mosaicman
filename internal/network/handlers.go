package network

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mohammadne/mosaicman/internal"
	"github.com/mohammadne/mosaicman/internal/models"
	"github.com/mohammadne/mosaicman/internal/mosaic"
	"github.com/mohammadne/mosaicman/pkg/logger"
	"github.com/mohammadne/mosaicman/pkg/utils"
)

type responseErr struct {
	Message string `json:"message"`
	Help    string `json:"help"`
}

var (
	missingImage = responseErr{
		Message: "image file is missed",
	}

	formatErr = responseErr{
		Message: "invalid image format has given",
		Help:    "right now we only support jpg format",
	}

	persisErr = responseErr{
		Message: "error persisting file",
		Help:    "please try later",
	}
)

func (server *server) upload(c echo.Context) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, missingImage)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	parts := strings.Split(fileHeader.Filename, ".")
	if parts[len(parts)-1] != "jpg" {
		return c.JSON(http.StatusBadRequest, formatErr)
	}

	metadata := &models.Metadata{
		IP:   c.Request().RemoteAddr,
		UUID: utils.NewUUID(),
	}

	if err = server.storage.Persist(context.TODO(), file, metadata); err != nil {
		return c.JSON(http.StatusBadRequest, persisErr)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":    "Image uploaded successfully",
		"image_uuid": metadata.UUID,
	})
}

var (
	invalidDataErr = responseErr{
		Message: "image file is missed",
	}

	imageNotFoundErr = responseErr{
		Message: "image file not found",
		Help:    "please upload your new image and try again",
	}

	processErr = responseErr{
		Message: "an internal server error occured",
		Help:    "please try later",
	}
)

func (server *server) process(c echo.Context) error {
	uuid := c.Param("image-uuid")
	if uuid == "" {
		return c.String(http.StatusBadRequest, "image uuid is missing")
	}

	requestData := new(models.Options)
	if err := c.Bind(requestData); err != nil {
		return c.JSON(http.StatusBadRequest, invalidDataErr)
	}

	metadata := &models.Metadata{IP: c.Request().RemoteAddr, UUID: uuid}
	original, err := server.storage.Retrieve(context.TODO(), metadata)
	if err != nil {
		server.logger.Error("error in retrieving original", logger.Error(err))
		return c.JSON(http.StatusBadRequest, imageNotFoundErr)
	}
	defer original.Close()

	err = mosaic.Process(original, uuid, requestData, server.tiles)
	if err != nil {
		server.logger.Error("error in processing mosaic", logger.Error(err))
		return c.JSON(http.StatusInternalServerError, processErr)
	}

	return c.String(http.StatusCreated, "mosaic image has been created")
}

func (server *server) get(c echo.Context) error {
	uuid := c.Param("image-uuid")
	if uuid == "" {
		return c.String(http.StatusBadRequest, "image uuid is missing")
	}

	metadata := &models.Metadata{
		IP:   c.Request().RemoteAddr,
		UUID: uuid,
	}

	err := server.storage.Validate(context.TODO(), metadata)
	if err != nil {
		server.logger.Error("error in retrieving original", logger.Error(err))
		return c.JSON(http.StatusBadRequest, imageNotFoundErr)
	}

	path := fmt.Sprintf("%s/%s.jpg", internal.MosaicsDir, uuid)
	return c.Attachment(path, "mosaic-result.jpg")
}
