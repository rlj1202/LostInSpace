package lostinspace

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
)

// A sector is a set of chunks.
// One sector will be saved into one file.
type Sector struct {
	Chunks [SECTOR_WIDTH * SECTOR_HEIGHT]*Chunk

	coord WorldSectorCoord
}

func NewSector(coord WorldSectorCoord) *Sector {
	sector := &Sector{
		coord: coord,
	}

	return sector
}

func (sector *Sector) Set(chunk *Chunk) {
	if !chunk.coord.Valid() {
		return
	}
	sector.Chunks[chunk.coord.X+chunk.coord.Y*SECTOR_WIDTH] = chunk
}

func (sector *Sector) At(coord ChunkCoord) *Chunk {
	if !coord.Valid() {
		return nil
	}
	return sector.Chunks[coord.X+coord.Y*SECTOR_WIDTH]
}

func SaveSector(sector *Sector) { // TODO
	file, err := os.Create(fmt.Sprintf("sector_%d_%d.gob", sector.coord.X, sector.coord.Y))
	if err != nil {
		panic(err)
	}

	enc := gob.NewEncoder(file)
	enc.Encode(sector)

	file.Close()

	jsonFile, err := os.Create("sector_0_0.json")
	jsonEnc := json.NewEncoder(jsonFile)
	jsonEnc.SetIndent("", "\t")
	jsonEnc.Encode(sector)
	jsonFile.Close()
}
