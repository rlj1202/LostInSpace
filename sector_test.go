package lostinspace_test

import (
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
					chunk.Set(lostinspace.NewBlock(lostinspace.BlockCoord{blockX, blockY}, "stone"))
				}
			}
			sector.Set(chunk)
		}
	}

	lostinspace.SaveSector(sector)

	t.Logf("Save file\n")
}
