package storage

import (
	"context"
	"fmt"
	"io"

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
	name := fmt.Sprintf("%s.jpg", uuid)
	path := fmt.Sprintf("%s/%s", s.savePath, name)

	destination, err := utils.CreateFile(path)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err = io.Copy(destination, file); err != nil {
		return err
	}

	return nil
}

func (storage *storage) Retrieve(ctx context.Context, md *models.Metadata) (interface{}, error) {
	return nil, nil
}
