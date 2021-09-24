package tiles

import "github.com/mohammadne/mosaicman/pkg/logger"

type Tiles interface {
	CloneTilesDB() Database
}

type tiles struct {
	logger   logger.Logger
	database Database
}

func New(lg logger.Logger) (Tiles, error) {
	t := &tiles{logger: lg}

	var err error
	if t.database, err = loadTilesDB(); err != nil {
		return nil, err
	}

	return t, nil
}
