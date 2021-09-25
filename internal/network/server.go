package network

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadne/mosaicman/internal/storage"
	"github.com/mohammadne/mosaicman/internal/tiles"
	"github.com/mohammadne/mosaicman/pkg/logger"
)

type server struct {
	address string
	storage storage.Storage
	tiles   tiles.Tiles
	logger  logger.Logger
	echo    *echo.Echo
}

func New(address string, s storage.Storage, t tiles.Tiles, lg logger.Logger) *server {
	server := &server{address: address, storage: s, tiles: t, logger: lg}

	server.echo = echo.New()
	server.echo.HideBanner = true

	server.echo.POST("/upload", server.upload)
	server.echo.POST("/process/:image-uuid", server.process)
	server.echo.GET("/mosaics/:image-uuid", server.get)

	return server
}

func (server *server) Serve() error {
	server.logger.Info("starting server", logger.String("on", server.address))
	if err := server.echo.Start(server.address); err != nil {
		return err
	}

	return nil
}
