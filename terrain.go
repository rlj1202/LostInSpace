package main

import (
	"encoding/gob"
	"math/rand"
	"os"
)

// Contains chunks which contain blocks.
type Terrain struct {
	Seed   int64
	Chunks map[Coord]*Chunk
}

// Represent coordinates.
type Coord struct {
	X, Y int64
}

func (terrain *Terrain) GetBlock(x, y int64) *Block {
	chunkCoord := Coord{x / CHUNK_WIDTH, y / CHUNK_HEIGHT}

	return terrain.Chunks[chunkCoord].At(x&0xf, y&0xf)
}

func (terrain *Terrain) SetBlock(block *Block) {
	x := int64(block.X)
	y := int64(block.Y)
	chunkCoord := Coord{x / CHUNK_WIDTH, y / CHUNK_HEIGHT}

	terrain.Chunks[chunkCoord].Set(block)
}

func (terrain *Terrain) GenerateChunk(coord Coord) *Chunk {
	random := rand.New(rand.NewSource(terrain.Seed))
	perm := [256]int{}
	copy(perm[:], random.Perm(256))
	chunk := NewChunk()
	chunk.X = coord.X
	chunk.Y = coord.Y

	for y := int64(0); y < CHUNK_HEIGHT; y++ {
		for x := int64(0); x < CHUNK_WIDTH; x++ {
			block := new(Block)
			block.X = x
			block.Y = y
			block.BlockType = ""

			blockX := float64(coord.X*CHUNK_WIDTH + x)
			blockY := float64(coord.Y*CHUNK_HEIGHT + y)

			noise := PerlinNoiseImproved(perm,
				blockX/8.0,
				blockY/8.0,
				0)

			if noise > 0.75 {
				block.BlockType = "test1"
			} else if noise > 0.7 {
				block.BlockType = "test2"
			} else if noise > 0.65 {
				block.BlockType = "stone"
			}

			chunk.Set(block)
		}
	}

	return chunk
}

func SaveTerrain(terrain *Terrain) {
	file, err := os.Create("terrain.gob")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	encoder := gob.NewEncoder(file)
	encoder.Encode(terrain)
}

func LoadTerrain() (*Terrain, error) {
	file, err := os.Open("terrain.gob")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	terrain := new(Terrain)

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(terrain)

	return terrain, err
}
