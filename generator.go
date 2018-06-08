package lostinspace

func GenerateSector(seed *Seed, sectorCoord WorldSectorCoord) *Sector {
	sector := NewSector(sectorCoord)

	for y := uint8(0); y < SECTOR_HEIGHT; y++ {
		for x := uint8(0); x < SECTOR_WIDTH; x++ {
			chunkCoord := ChunkCoord{x, y}
			worldChunkCoord := CombineWorldChunkCoord(sectorCoord, chunkCoord)

			chunk := GenerateChunk(seed, worldChunkCoord)
			sector.Set(chunk)
		}
	}

	return sector
}

func GenerateChunk(seed *Seed, worldChunkCoord WorldChunkCoord) *Chunk {
	perm := seed.perm

	sectorCoord, chunkCoord := worldChunkCoord.Parse()

	chunk := NewChunk(chunkCoord)

	for y := uint8(0); y < CHUNK_HEIGHT; y++ {
		for x := uint8(0); x < CHUNK_WIDTH; x++ {
			blockCoord := BlockCoord{x, y}
			block := NewBlock(blockCoord, BLOCK_TYPE_VOID, 0)

			worldCoord := CombineWorldBlockCoord(sectorCoord, chunkCoord, blockCoord)

			noise := PerlinNoiseImproved(perm,
				float64(worldCoord.X)/32.0,
				float64(worldCoord.Y)/32.0,
				0)

			if noise > 0.67 {
				noise *= PerlinNoiseImproved(
					perm,
					float64(worldCoord.X)/8.0,
					float64(worldCoord.Y)/8.0,
					0,
				)

				if noise > 0.5 {
					block.BlockType = "test1"
				} else if noise > 0.45 {
					block.BlockType = "test2"
				} else {
					block.BlockType = "stone"
				}
			}

			chunk.Set(block)
		}
	}

	return chunk
}
