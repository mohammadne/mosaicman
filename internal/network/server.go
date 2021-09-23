package network

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadne/mosaicman/pkg/logger"
	"github.com/mohammadne/mosaicman/pkg/utils"
)

type server struct {
	address  string
	savePath string
	logger   logger.Logger
	echo     *echo.Echo
}

func New(address string, savePath string, lg logger.Logger) *server {
	server := &server{address: address, savePath: savePath, logger: lg}

	server.echo = echo.New()
	server.echo.HideBanner = true

	server.echo.POST("/upload", server.upload)
	server.echo.POST("/generate", server.generate)

	return server
}

func (server *server) Serve() error {
	err := utils.CreateDirIfMissed(server.savePath)
	if err != nil {
		return err
	}

	server.logger.Info("starting server", logger.String("on", server.address))
	if err := server.echo.Start(server.address); err != nil {
		return err
	}

	return nil
}
