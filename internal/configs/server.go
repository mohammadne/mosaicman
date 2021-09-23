package configs

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mohammadne/mosaicman/pkg/logger"
)

type server struct {
	Address string
	Logger  *logger.Config
}

func Server(env string) *server {
	config := &server{}

	switch env {
	case "prod":
		config.loadProd()
	default:
		config.loadDev()
	}

	return config
}

func (config *server) loadProd() {
	config.Logger = &logger.Config{}

	// process
	envconfig.MustProcess("", config)
	envconfig.MustProcess("logger", config.Logger)

}

func (config *server) loadDev() {
	config.Address = "localhost:8080"

	config.Logger = &logger.Config{
		Development:      true,
		EnableCaller:     true,
		EnableStacktrace: false,
		Encoding:         "console",
		Level:            "warn",
	}
}
