package mosaic

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"

	"github.com/mohammadne/mosaicman/internal"
	"github.com/mohammadne/mosaicman/internal/models"
	"github.com/mohammadne/mosaicman/internal/tiles"
	"github.com/mohammadne/mosaicman/pkg/utils"
)

func Process(file io.Reader, uuid string, opts models.Options, t tiles.Tiles) error {
	original, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	db := t.CloneTilesDB()
	bounds := original.Bounds()

	// fan-out
	part1 := fanOut(original, &db, opts.TileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
	part2 := fanOut(original, &db, opts.TileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
	part3 := fanOut(original, &db, opts.TileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
	part4 := fanOut(original, &db, opts.TileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)

	// fan-in
	combine := fanIn(bounds, part1, part2, part3, part4)

	mosaic, err := jpeg.Decode(<-combine)
	if err != nil {
		return err
	}

	mask := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	res := backgroundBellowForeground(mask, original, mosaic)

	return persist(res, uuid)
}

func persist(image image.Image, uuid string) error {
	path := fmt.Sprintf("%s/%s.jpg", internal.MosaicsDir, uuid)
	destination, err := utils.CreateFile(path)
	if err != nil {
		return err
	}
	defer destination.Close()

	err = jpeg.Encode(destination, image, nil)
	if err != nil {
		return err
	}

	return nil
}
