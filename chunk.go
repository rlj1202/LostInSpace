package lostinspace

type Chunk struct {
	Blocks [CHUNK_WIDTH * CHUNK_HEIGHT]*Block

	coord ChunkCoord

	mesh *Mesh
	body *Body

	aabb *AABB
}

func NewChunk(coord ChunkCoord) *Chunk {
	chunk := new(Chunk)
	chunk.coord = coord
	for y := 0; y < CHUNK_HEIGHT; y++ {
		for x := 0; x < CHUNK_WIDTH; x++ {
			chunk.Set(NewBlock(BlockCoord{uint8(x), uint8(y)}, BLOCK_TYPE_VOID))
		}
	}

	return chunk
}

func (chunk *Chunk) Set(block *Block) {
	if !block.coord.Valid() {
		return
	}

	chunk.Blocks[blockIndex(block.coord)] = block
}

func (chunk *Chunk) At(coord BlockCoord) *Block {
	if !coord.Valid() {
		return nil
	}

	return chunk.Blocks[blockIndex(coord)]
}

func (chunk *Chunk) ForEach(f func(*Block)) {
	for _, block := range chunk.Blocks {
		f(block)
	}
}

// Deallocate vao, vbo, ebo and b2body.
func (chunk *Chunk) Destroy() {
	chunk.mesh.Destroy()
	chunk.body.Destroy()
	chunk.mesh = nil
	chunk.body = nil
}

func (chunk *Chunk) GetAABB() *AABB {
	if chunk.aabb == nil {
		x, y := float64(chunk.coord.X*16+8), float64(chunk.coord.Y*16+8)
		chunk.aabb = &AABB{
			Center:  Vec2{x, y},
			HWidth:  8,
			HHeight: 8,
		}
	}

	return chunk.aabb
}

func (chunk *Chunk) Bake(world *World, dic *BlockTypeDictionary, coord WorldChunkCoord) {
	if chunk.mesh == nil {
		chunk.mesh = NewMesh(nil, nil, nil, nil)
	}
	BakeBlockStorageMesh(chunk.mesh, chunk, dic)

	if chunk.body == nil {
		chunk.body = world.CreateBody(false)
		chunk.body.SetPosition(
			float64(coord.X*CHUNK_WIDTH),
			float64(coord.Y*CHUNK_HEIGHT),
		)
	}
	BakeBlockStorageBody(chunk.body, chunk, dic)
}

func blockIndex(coord BlockCoord) int {
	return int(coord.X + coord.Y*CHUNK_WIDTH)
}
