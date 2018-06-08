package lostinspace

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

// Sector manager loads sectors from file.
// If there is no file, manager will generate one.
func sectorManager(game *Game) {
	for {
		select {
		case <-game.quit:
			return
		default:
			x, y := game.player.GetPosition()
			worldBlockCoord := WorldBlockCoord{
				int64(math.Floor(x)),
				int64(math.Floor(y)),
			}
			curSectorCoord, _, _ := worldBlockCoord.Parse()

			// load sectors and request baking
			sectorCoordsToLoad := []WorldSectorCoord{
				curSectorCoord,
				curSectorCoord.Left(),
				curSectorCoord.Left().Up(),
				curSectorCoord.Left().Down(),
				curSectorCoord.Right(),
				curSectorCoord.Right().Up(),
				curSectorCoord.Right().Down(),
				curSectorCoord.Up(),
				curSectorCoord.Down(),
			}
			for _, coord := range sectorCoordsToLoad {
				sector := game.terrain.GetSector(coord)
				if sector != nil {
					continue
				}

				sector = LoadSector(coord)
				if sector == nil {
					sector = GenerateSector(game.terrain.Seed, coord)
				}

				for _, chunk := range sector.Chunks {
					worldChunkCoord := CombineWorldChunkCoord(coord, chunk.coord)
					chunk.Bake(game.world, game.dic, worldChunkCoord)

					game.bakeChunkQueue <- chunk
				}

				game.terrain.SetSector(sector)

				log.Printf("Load %v\n", coord)
			}

			// unload sectors and request destroying
			loadedSectorCoords := make(map[WorldSectorCoord]bool)
			for _, sector := range game.terrain.Sectors {
				loadedSectorCoords[sector.coord] = true
			}
			for _, sectorCoord := range sectorCoordsToLoad {
				delete(loadedSectorCoords, sectorCoord)
			}

			for sectorCoord := range loadedSectorCoords {
				sector := game.terrain.GetSector(sectorCoord)
				if sector == nil {
					continue
				}

				game.terrain.DeleteSector(sectorCoord)

				SaveSector(sector)

				for _, chunk := range sector.Chunks {
					game.destroyChunkQueue <- chunk
				}

				log.Printf("Unload %v\n", sectorCoord)
			}

			time.Sleep(time.Second / 2)
		}
	}
}

func sectorFileName(coord WorldSectorCoord) string {
	return fmt.Sprintf("sector_%d_%d.gob", coord.X, coord.Y)
}

func LoadSector(coord WorldSectorCoord) *Sector {
	file, err := os.Open(sectorFileName(coord))
	if err != nil {
		return nil
	}
	defer file.Close()

	sector := NewSector(coord)

	dec := gob.NewDecoder(file)
	dec.Decode(sector)

	sector.coord = coord
	for i, chunk := range sector.Chunks {
		chunk.coord = ChunkCoord{
			uint8(i % SECTOR_WIDTH),
			uint8(i / SECTOR_WIDTH),
		}

		for j, block := range chunk.Blocks {
			block.coord = BlockCoord{
				uint8(j % CHUNK_WIDTH),
				uint8(j / CHUNK_WIDTH),
			}
		}
	}

	return sector
}

func SaveSector(sector *Sector) {
	file, err := os.Create(sectorFileName(sector.coord))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	enc.Encode(sector)
}
