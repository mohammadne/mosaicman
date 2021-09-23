package network

import (
	"github.com/labstack/echo/v4"
)

type server struct {
	address string
	echo    *echo.Echo
}

func New(address string) *server {
	server := &server{address: address}

	server.echo = echo.New()
	server.echo.HideBanner = true
	server.setupRoutes()

	return server
}

func (server *server) setupRoutes() {}

func (server *server) Serve() {}
