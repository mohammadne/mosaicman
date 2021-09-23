package server

import (
	"github.com/mohammadne/mosaicman/internal/configs"
	"github.com/mohammadne/mosaicman/internal/network"
	"github.com/mohammadne/mosaicman/internal/storage"
	"github.com/mohammadne/mosaicman/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	use   = "server"
	short = "run server"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{Use: use, Short: short, Run: main}

	envFlag := "set config environment, default is dev"
	cmd.Flags().StringP("env", "e", "", envFlag)

	return cmd
}

func main(cmd *cobra.Command, _ []string) {
	env := cmd.Flag("env").Value.String()
	configs := configs.Server(env)

	lg := logger.NewZap(configs.Logger)

	storage, err := storage.New(nil, configs.SavePath, lg)
	if err != nil {
		lg.Fatal("error creating redis storage", logger.Error(err))
	}

	server := network.New(configs.Address, storage, lg)
	if err := server.Serve(); err != nil {
		lg.Fatal("starting server failed", logger.Error(err))
	}
}
