package lostinspace

// Terrian is set of chunks.
type Terrain struct {
	Chunks map[ChunkCoord]*Chunk
}

func NewTerrain() *Terrain {
	terrain := &Terrain{
		Chunks: make(map[ChunkCoord]*Chunk),
	}

	return terrain
}

func (terrain *Terrain) SetChunk(chunk *Chunk) {
	terrain.Chunks[chunk.ChunkCoord] = chunk
}

func (terrain *Terrain) GetChunk(coord ChunkCoord) *Chunk {
	chunk := terrain.Chunks[coord]

	return chunk
}

// Set block to given world coord.
// Block coord which a block has will be ignored by world coord.
func (terrain *Terrain) SetBlock(coord WorldCoord, block *Block) {
	chunkCoord, blockCoord := coord.Parse()

	chunk, exist := terrain.Chunks[chunkCoord]
	if !exist {
		// do something
		return
	}

	block.BlockCoord = blockCoord

	chunk.Set(block)
}

// Get block at given world coord.
// It will return nil if there is no corresponding chunk.
func (terrain *Terrain) At(coord WorldCoord) *Block {
	chunkCoord, blockCoord := coord.Parse()

	chunk, exist := terrain.Chunks[chunkCoord]
	if !exist {
		return nil
	}

	return chunk.At(blockCoord)
}
