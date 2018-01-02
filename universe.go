package lostinspace

import (
	"os"
)

// Yes, universe. This is UNIVERSE.
// Contains whole things of the world.
//
//  universe/          # a directory which contains all infos about an universe.
//      data.gob       # blah blah
//      somedata.gob
//      blahblah.gob
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
	dir *os.File

	*World
	*Terrain
	player   *Player
	entities []*BlockEntity
}

func NewUniverse(dir *os.File, playerTexFile *os.File) *Universe {
	playerTex := NewTexture2D(playerTexFile)

	universe := new(Universe)
	universe.dir = dir
	universe.World = NewWorld()
	universe.player = NewPlayer(universe.World, playerTex)
	universe.Terrain = NewTerrain()
	universe.entities = make([]*BlockEntity, 0)

	return universe
}

func (universe *Universe) CreateBlockEntity() *BlockEntity {
	entity := NewBlockEntity(universe.World)

	return entity
}

// Save all informations to dir.
func (universe *Universe) Save() {
}

// Load all informations from dir.
func (universe *Universe) Load() {
}
