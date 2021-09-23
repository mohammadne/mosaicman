package storage

import (
	"context"
	"io"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mohammadne/mosaicman/internal/models"
	"github.com/mohammadne/mosaicman/pkg/logger"
)

type Storage interface {
	Persist(context.Context, io.Reader, *models.Metadata) error
	Retrieve(context.Context, *models.Metadata) (string, error)
}

type storage struct {
	logger  logger.Logger
	expires time.Duration
	cmd     redis.Cmdable
}

func New(cfg *Config, lg logger.Logger) (Storage, error) {
	expires := time.Duration(cfg.ExpirationTime)
	s := &storage{logger: lg, expires: time.Second * expires}

	if Mode(cfg.Mode) == Cluster {
		s.cmd = newClusterRedis(cfg)
	} else {
		s.cmd = newSingleRedis(cfg)
	}

	if err := s.cmd.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}

	return s, nil
}
