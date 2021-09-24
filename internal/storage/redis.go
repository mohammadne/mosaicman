package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/gomodule/redigo/redis"
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

	connection, err := s.newConnection(ctx)
	if err != nil {
		return err
	}
	defer connection.Close()

	if _, err := connection.Do("SET", md.UUID, md.IP, "EX", s.config.Expiration); err != nil {
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
	connection, err := s.newConnection(ctx)
	if err != nil {
		return "", err
	}
	defer connection.Close()

	value, err := redis.String(connection.Do("GET", md.UUID))
	if err != nil {
		if err == redis.ErrNil {
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
