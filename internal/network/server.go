package network

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadne/mosaicman/internal/storage"
	"github.com/mohammadne/mosaicman/pkg/logger"
)

type server struct {
	address string
	storage storage.Storage
	logger  logger.Logger
	echo    *echo.Echo
}

func New(address string, storage storage.Storage, lg logger.Logger) *server {
	server := &server{address: address, storage: storage, logger: lg}

	server.echo = echo.New()
	server.echo.HideBanner = true

	server.echo.POST("/upload", server.upload)
	server.echo.POST("/process", server.process)

	return server
}

func (server *server) Serve() error {
	server.logger.Info("starting server", logger.String("on", server.address))
	if err := server.echo.Start(server.address); err != nil {
		return err
	}

	return nil
}
