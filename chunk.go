package main

import (
	"github.com/ByteArena/box2d"
)

const (
	CHUNK_WIDTH  = 16
	CHUNK_HEIGHT = 16
)

// Chunk contains 16 * 16 blocks to control terrain easily.
type Chunk struct {
	Coord

	Blocks []*Block
	Width  int
	Height int

	*BlockContainerObject
}

func NewChunk() *Chunk {
	return &Chunk{
		Blocks: make([]*Block, CHUNK_WIDTH*CHUNK_HEIGHT),
		Width:  CHUNK_WIDTH,
		Height: CHUNK_HEIGHT,
	}
}

// Set block at given coordinate.
// block's coordinates have to be in 0 ~ 15
func (chunk *Chunk) Set(block *Block) {
	x := block.X
	y := block.Y
	if x < 0 || CHUNK_WIDTH <= x || y < 0 || CHUNK_HEIGHT <= y {
		// panic
		return
	}
	chunk.Blocks[x+y*CHUNK_WIDTH] = block
}

// Get block at given coordinate.
func (chunk *Chunk) At(x, y int64) *Block {
	if x < 0 || CHUNK_WIDTH <= x || y < 0 || CHUNK_HEIGHT <= y {
		// panic
		return nil
	}
	return chunk.Blocks[x+y*CHUNK_WIDTH]
}

func (chunk *Chunk) ForEach(f func(*Block)) {
	for _, v := range chunk.Blocks {
		f(v)
	}
}

// Generate initial block meshs and physics box
func (chunk *Chunk) bake(game *Game) {
	object := NewBlockContainerObject(game, chunk, BLOCK_CONTAINER_STATIC)
	object.body.SetTransform(box2d.MakeB2Vec2(
		float64(chunk.Coord.X*CHUNK_WIDTH),
		float64(chunk.Coord.Y*CHUNK_HEIGHT),
	), 0)

	chunk.BlockContainerObject = object
}
