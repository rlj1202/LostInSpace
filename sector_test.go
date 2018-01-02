package lostinspace_test

import (
	"fmt"
	"testing"

	"github.com/rlj1202/LostInSpace"
)

func TestSector(t *testing.T) {
	sector := lostinspace.NewSector(lostinspace.WorldSectorCoord{0, 0})

	for chunkY := uint8(0); chunkY < 16; chunkY++ {
		for chunkX := uint8(0); chunkX < 16; chunkX++ {
			chunk := lostinspace.NewChunk(lostinspace.ChunkCoord{chunkX, chunkY})

			for blockY := uint8(0); blockY < 16; blockY++ {
				for blockX := uint8(0); blockX < 16; blockX++ {
					chunk.Set(lostinspace.NewBlock(
						lostinspace.BlockCoord{blockX, blockY},
						lostinspace.BlockType(fmt.Sprintf("stone_%d_%d", blockX, blockY)),
					))
				}
			}
			sector.Set(chunk)
		}
	}

	lostinspace.SaveSector(sector)
	t.Logf("Saved file\n")

	newSector := lostinspace.LoadSector(lostinspace.WorldSectorCoord{0, 0})
	t.Logf("Loaded file\n")

	for i, block := range newSector.Chunks[0].Blocks {
		t.Logf("[%3d]: %v\n", i, block.BlockType)
	}
}
