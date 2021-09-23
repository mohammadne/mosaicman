package configs

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mohammadne/mosaicman/pkg/logger"
)

type server struct {
	Address  string
	SavePath string `split_words:"true"`
	Logger   *logger.Config
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

	// process
	envconfig.MustProcess("", configs)
	envconfig.MustProcess("logger", configs.Logger)

}

func (configs *server) loadDev() {
	configs.Address = "localhost:8080"

	configs.SavePath = "assets/uploads"

	configs.Logger = &logger.Config{
		Development:      true,
		EnableCaller:     true,
		EnableStacktrace: false,
		Encoding:         "console",
		Level:            "warn",
	}
}
