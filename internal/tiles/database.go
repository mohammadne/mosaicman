package tiles

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Database struct {
	Mutex *sync.Mutex
	Store DataModel
}

type DataModel map[string][3]float64

// loadTilesDB populates tiles database in memory
func loadTilesDB() (Database, error) {
	fmt.Println("Start populating tiles db ...")

	db := make(DataModel)
	files, _ := ioutil.ReadDir("static/tiles/cats")
	for _, f := range files {
		name := filepath.Join("static/tiles/cats", f.Name())
		file, err := os.Open(name)
		if err == nil {
			img, _, err := image.Decode(file)
			if err == nil {
				db[name] = averageColor(img)
			} else {
				fmt.Println("error in populating tiles db:", err, name)
			}
		} else {
			fmt.Println("cannot open file", name, "when populating tiles db:", err)
		}
		file.Close()
	}

	fmt.Println("Finished populating tiles db.")

	return Database{
		Store: db,
		Mutex: &sync.Mutex{},
	}, nil
}

func (tiles *tiles) CloneTilesDB() Database {
	db := make(map[string][3]float64)
	for k, v := range tiles.database.Store {
		db[k] = v
	}

	return Database{
		Store: db,
		Mutex: &sync.Mutex{},
	}
}
