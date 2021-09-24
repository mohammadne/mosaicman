package mosaic

import (
	"image"
	"image/color"
	"math"

	"github.com/mohammadne/mosaicman/internal/tiles"
)

// resize an image by its ratio e.g. ratio 2 means reduce the size by 1/2, 10 means reduce the size by 1/10
func resize(in image.Image, newWidth int) image.NRGBA {
	bounds := in.Bounds()
	width := bounds.Dx()
	ratio := width / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.X/ratio, bounds.Max.X/ratio, bounds.Max.Y/ratio))
	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+ratio, i+1 {
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	return *out
}

func nearest(db *tiles.Database, target [3]float64) string {
	var filename string
	db.Mutex.Lock()
	smallest := 1000000.0
	for k, v := range db.Store {
		dist := distance(target, v)
		if dist < smallest {
			filename, smallest = k, dist
		}
	}
	delete(db.Store, filename)
	db.Mutex.Unlock()
	return filename
}

// find the Eucleadian distance between 2 points
func distance(p1 [3]float64, p2 [3]float64) float64 {
	return math.Sqrt(sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2]))
}

// find the square
func sq(n float64) float64 {
	return n * n
}
