package mosaic

import (
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/mohammadne/mosaicman/internal/models"
	"github.com/mohammadne/mosaicman/internal/tiles"
)

func Process(org io.Reader, opts models.Options, t tiles.Tiles) (string, error) {
	image, _, err := image.Decode(org)
	if err != nil {
		return "", err
	}

	db := t.CloneTilesDB()
	bounds := image.Bounds()

	// fan-out
	part1 := fanOut(image, &db, opts.TileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
	part2 := fanOut(image, &db, opts.TileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
	part3 := fanOut(image, &db, opts.TileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
	part4 := fanOut(image, &db, opts.TileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)

	// fan-in
	combine := fanIn(bounds, part1, part2, part3, part4)

	jpgI, err := jpeg.Decode(<-combine)
	if err != nil {
		return "", err
	}

	out, _ := os.Create("./img.jpeg")
	defer out.Close()

	err = jpeg.Encode(out, jpgI, nil)
	if err != nil {
		return "", err
	}

	return "./img.jpeg", nil
}
