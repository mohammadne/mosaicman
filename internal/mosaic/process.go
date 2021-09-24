package mosaic

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"

	"github.com/mohammadne/mosaicman/internal/models"
)

func Process(original io.Reader, options models.Options) (interface{}, error) {
	image, _, err := image.Decode(original)
	if err != nil {
		return "", err
	}

	db := cloneTilesDB()
	bounds := image.Bounds()

	// fan-out
	part1 := fanOut(image, &db, options.TileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
	part2 := fanOut(image, &db, options.TileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
	part3 := fanOut(image, &db, options.TileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
	part4 := fanOut(image, &db, options.TileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)

	// fan-in
	combine := fanIn(bounds, part1, part2, part3, part4)

	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, image, nil)

	return map[string]string{
		"original": base64.StdEncoding.EncodeToString(buffer.Bytes()),
		"mosaic":   <-combine,
	}, nil
}
