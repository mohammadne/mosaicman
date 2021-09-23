package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-redis/redis/v8"
	"github.com/mohammadne/mosaicman/internal"
	"github.com/mohammadne/mosaicman/internal/models"
	"github.com/mohammadne/mosaicman/pkg/logger"
	"github.com/mohammadne/mosaicman/pkg/utils"
)

func (s *storage) Persist(ctx context.Context, file io.Reader, md *models.Metadata) error {
	if err := s.persistFile(ctx, file, md.UUID); err != nil {
		s.logger.Error("error saving file", logger.Error(err))
		return err
	}

	if err := s.cmd.Set(ctx, md.UUID, md.IP, s.expires).Err(); err != nil {
		s.logger.Error("error saving metadata", logger.Error(err))
		return err
	}

	return nil
}

func (s *storage) persistFile(ctx context.Context, file io.Reader, uuid string) error {
	destination, err := utils.CreateFile(s.getPath(uuid))
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err = io.Copy(destination, file); err != nil {
		return err
	}

	return nil
}

func (s *storage) Retrieve(ctx context.Context, md *models.Metadata) (string, error) {
	value, err := s.cmd.Get(ctx, md.UUID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("no matching record found in redis database")
		}

		s.logger.Error("error getting from redis", logger.Error(err))
		return "", err
	}

	if value != md.IP {
		err = errors.New("requster IP doesn't match with requested image IP")
		s.logger.Error("error getting from redis", logger.Error(err))
		return "", err
	}

	return s.getPath(md.UUID), nil
}

func (s *storage) getPath(uuid string) string {
	name := fmt.Sprintf("%s.jpg", uuid)
	return fmt.Sprintf("%s/%s", internal.OriginalsDir, name)
}
