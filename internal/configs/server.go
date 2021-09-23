package configs

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mohammadne/mosaicman/internal/storage"
	"github.com/mohammadne/mosaicman/pkg/logger"
)

type server struct {
	Address string
	Logger  *logger.Config
	Storage *storage.Config
}

func Server(env string) *server {
	configs := &server{}

	switch env {
	case "prod":
		configs.loadProd()
	default:
		configs.loadDev()
	}

	return configs
}

func (configs *server) loadProd() {
	configs.Logger = &logger.Config{}
	configs.Storage = &storage.Config{}

	// process
	envconfig.MustProcess("", configs)
	envconfig.MustProcess("logger", configs.Logger)
	envconfig.MustProcess("storage", configs.Storage)

}

func (configs *server) loadDev() {
	configs.Address = "localhost:8080"

	configs.Logger = &logger.Config{
		Development:      true,
		EnableCaller:     true,
		EnableStacktrace: false,
		Encoding:         "console",
		Level:            "warn",
	}

	configs.Storage = &storage.Config{
		Mode:           0,
		URL:            "localhost:6379",
		ExpirationTime: 10,
	}
}
