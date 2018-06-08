package lostinspace

import (
	"os"
	"time"
)

// Yes, universe. This is UNIVERSE.
// Contains whole things of the world.
//
//  universe/              # A directory which contains all informations about an universe.
//      sectors/           # Each sector is saved into one file.
//          sector_x_y.gob
//          sector_x_y.gob
//          ...
//
// data structure
//
//  Universe: {
//      Terrain: {
//          Sectors: map[WorldSectorCoord]Sector{
//              {
//                  Chunks: [8*8]Chunk{
//                      {
//                          Blocks: [16*16]Block{
//                              {
//                                  BlockType: string,
//                              },
//                              ...
//                          },
//                      },
//                      ...
//                  },
//              },
//              ...
//          },
//      },
//      Entities: []Entity{
//          {
//              Blocks: []Block{},
//          },
//      }
//      Player: {
//          Inventory: {
//              Items: []Item{},
//          },
//      },
//  }
type Universe struct {
	*Terrain

	dir *os.File

	world    *World
	player   *Player
	entities []*BlockEntity
}

func NewUniverse(dir *os.File, playerTexFile *os.File) *Universe {
	playerTex := NewTexture2D(playerTexFile)

	universe := new(Universe)
	universe.dir = dir
	universe.world = NewWorld()
	universe.player = NewPlayer(universe.world, playerTex)
	universe.Terrain = NewTerrain()
	universe.entities = make([]*BlockEntity, 0)

	return universe
}

func (universe *Universe) CreateBlockEntity() *BlockEntity {
	entity := NewBlockEntity(universe.world)

	return entity
}

func (universe *Universe) Update(dt time.Duration) {
	universe.world.Update(dt)
}

func (universe *Universe) SaveSector(worldSectorCoord WorldSectorCoord) {
}

func (universe *Universe) LoadSector(worldSectorCoord WorldSectorCoord) {
}
