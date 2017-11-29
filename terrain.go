package main

import (
	"encoding/gob"
	"math/rand"
	"os"
)

const (
	CHUNK_WIDTH  = 16
	CHUNK_HEIGHT = 16
)

// Contains chunks which contain blocks.
type Terrain struct {
	Seed   int64
	Chunks map[Coord]*Chunk
}

// Represent chunk coordinates.
type Coord struct {
	X, Y int64
}

// Chunk contains 16 * 16 blocks to control terrain easily.
type Chunk struct {
	Coord
	Blocks [CHUNK_WIDTH * CHUNK_HEIGHT]*Block
}

// Container for block informations such as block type, coordinates, etc.
type Block struct {
	BlockType
	X uint8
	Y uint8
	// Represents where block is facing.
	// 0 = right, 1 = up, 2 = left, 3 = bottom (counter-clock wise)
	FrontFace uint8
}

type BlockTypeDescriptor struct {
	BlockType
	// Block name which will be appeared to user.
	Name string
	// Each vertex must be in range between -0.5 and 0.5.
	CollisionVertices []Vec2
	TextureFile       *os.File
}

// Identifier to distinguish block type. It is string type.
type BlockType string

func (terrain *Terrain) GetBlock(x, y int64) *Block {
	chunkCoord := Coord{x / CHUNK_WIDTH, y / CHUNK_HEIGHT}

	return terrain.Chunks[chunkCoord].GetBlock(uint8(x&0xf), uint8(y&0xf))
}

func (terrain *Terrain) SetBlock(x, y int64, block *Block) {
	chunkCoord := Coord{x / CHUNK_WIDTH, y / CHUNK_HEIGHT}

	terrain.Chunks[chunkCoord].SetBlock(uint8(x&0xf), uint8(y&0xf), block)
}

func (chunk *Chunk) GetBlock(x, y uint8) *Block {
	return chunk.Blocks[x+y*CHUNK_WIDTH]
}

func (chunk *Chunk) SetBlock(x, y uint8, block *Block) {
	chunk.Blocks[x+y*CHUNK_WIDTH] = block
}

func (terrain *Terrain) GenerateChunk(coord Coord) *Chunk {
	random := rand.New(rand.NewSource(terrain.Seed))
	perm := [256]int{}
	copy(perm[:], random.Perm(256))
	chunk := new(Chunk)
	chunk.X = coord.X
	chunk.Y = coord.Y

	for y := int64(0); y < CHUNK_HEIGHT; y++ {
		for x := int64(0); x < CHUNK_WIDTH; x++ {
			block := new(Block)
			block.X = uint8(x)
			block.Y = uint8(y)
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

			chunk.SetBlock(uint8(x), uint8(y), block)
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
