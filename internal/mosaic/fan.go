package mosaic

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"os"
	"sync"

	"github.com/mohammadne/mosaicman/internal/tiles"
)

func fanIn(r image.Rectangle, c1, c2, c3, c4 <-chan image.Image) <-chan io.Reader {
	c := make(chan io.Reader)
	// start a goroutine
	go func() {
		var wg sync.WaitGroup
		newimage := image.NewNRGBA(r)
		copy := func(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
			draw.Draw(dst, r, src, sp, draw.Src)
			wg.Done()
		}
		wg.Add(4)
		var s1, s2, s3, s4 image.Image
		var ok1, ok2, ok3, ok4 bool
		for {
			select {
			case s1, ok1 = <-c1:
				go copy(newimage, s1.Bounds(), s1, image.Point{r.Min.X, r.Min.Y})
			case s2, ok2 = <-c2:
				go copy(newimage, s2.Bounds(), s2, image.Point{r.Max.X / 2, r.Min.Y})
			case s3, ok3 = <-c3:
				go copy(newimage, s3.Bounds(), s3, image.Point{r.Min.X, r.Max.Y / 2})
			case s4, ok4 = <-c4:
				go copy(newimage, s4.Bounds(), s4, image.Point{r.Max.X / 2, r.Max.Y / 2})
			}
			if ok1 && ok2 && ok3 && ok4 {
				break
			}
		}
		// wait till all copy goroutines are complete
		wg.Wait()
		buf2 := new(bytes.Buffer)
		jpeg.Encode(buf2, newimage, nil)
		c <- buf2
	}()

	return c
}

func fanOut(original image.Image, db *tiles.Database, tileSize, x1, y1, x2, y2 int) <-chan image.Image {
	c := make(chan image.Image)
	sp := image.Point{0, 0}
	go func() {
		newimage := image.NewNRGBA(image.Rect(x1, y1, x2, y2))
		for y := y1; y < y2; y = y + tileSize {
			for x := x1; x < x2; x = x + tileSize {
				r, g, b, _ := original.At(x, y).RGBA()
				color := [3]float64{float64(r), float64(g), float64(b)}
				nearest := nearest(db, color)
				file, err := os.Open(nearest)
				if err == nil {
					img, _, err := image.Decode(file)
					if err == nil {
						t := resize(img, tileSize)
						tile := t.SubImage(t.Bounds())
						tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
						draw.Draw(newimage, tileBounds, tile, sp, draw.Src)
					} else {
						fmt.Println("error in decoding nearest", err, nearest)
					}
				} else {
					fmt.Println("error opening file when creating mosaic:", nearest)
				}
				file.Close()
			}
		}
		c <- newimage.SubImage(newimage.Rect)
	}()

	return c
}
