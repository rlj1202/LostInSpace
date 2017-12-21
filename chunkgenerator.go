package lostinspace

import "math/rand"

func GenerateChunk(seed int64, chunkCoord ChunkCoord) *Chunk {
	random := rand.New(rand.NewSource(seed))
	perm := [256]int{}
	copy(perm[:], random.Perm(256))
	chunk := NewChunk()
	chunk.ChunkCoord = chunkCoord

	for y := uint8(0); y < CHUNK_HEIGHT; y++ {
		for x := uint8(0); x < CHUNK_WIDTH; x++ {
			block := NewBlock(BlockCoord{x, y}, BLOCK_TYPE_VOID)

			worldCoord := CombindWorldCoord(chunkCoord, block.BlockCoord)

			noise := PerlinNoiseImproved(perm,
				float64(worldCoord.X)/8.0,
				float64(worldCoord.Y)/8.0,
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
