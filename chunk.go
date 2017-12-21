package lostinspace

type Chunk struct {
	ChunkCoord
	Blocks [CHUNK_WIDTH * CHUNK_HEIGHT]*Block

	*Mesh
	*Body
}

func NewChunk() *Chunk {
	chunk := &Chunk{}
	for y := 0; y < CHUNK_HEIGHT; y++ {
		for x := 0; x < CHUNK_WIDTH; x++ {
			chunk.Set(NewBlock(BlockCoord{uint8(x), uint8(y)}, BLOCK_TYPE_VOID))
		}
	}

	return chunk
}

func (chunk *Chunk) Set(block *Block) {
	if !validBlockCoord(block.BlockCoord) {
		// TODO panic
		return
	}

	chunk.Blocks[blockIndex(block.BlockCoord)] = block
}

func (chunk *Chunk) At(coord BlockCoord) *Block {
	if !validBlockCoord(coord) {
		// TODO panic
		return nil
	}

	return chunk.Blocks[blockIndex(coord)]
}

func (chunk *Chunk) ForEach(f func(*Block)) {
	for _, block := range chunk.Blocks {
		f(block)
	}
}

func (chunk *Chunk) Bake(world *World, dic *BlockTypeDictionary) {
	positions, _, coords, indices := BakeBlockStorageMesh(chunk, dic)

	if chunk.Mesh == nil {
		chunk.Mesh = NewMesh(positions, nil, coords, indices)
	} else {
		chunk.Mesh.Set(positions, nil, coords, indices)
	}

	if chunk.Body == nil {
		chunk.Body = world.CreateBody(false)
		chunk.Body.SetPosition(
			float64(chunk.X*CHUNK_WIDTH), float64(chunk.Y*CHUNK_HEIGHT))
	}
	BakeBlockStorageBody(chunk.Body, chunk, dic)
}

func (chunk *Chunk) Destroy() {
	chunk.Mesh.Destroy()
	chunk.Body.Destroy()
	chunk.Mesh = nil
	chunk.Body = nil
}

func blockIndex(coord BlockCoord) int {
	return int(coord.X + coord.Y*CHUNK_WIDTH)
}

func validBlockCoord(coord BlockCoord) bool {
	return 0 <= coord.X && coord.X < CHUNK_WIDTH && 0 <= coord.Y && coord.Y < CHUNK_HEIGHT
}
