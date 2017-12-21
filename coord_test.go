package lostinspace_test

import (
	"github.com/rlj1202/LostInSpace/new/lostinspace"
	"testing"
)

var (
	worldCoords = []lostinspace.WorldCoord{
		{X: 0, Y: 1},
		{X: -1, Y: -1},
		{X: 16, Y: 16},
	}
	chunkCoords = []lostinspace.ChunkCoord{
		{X: 0, Y: 0},
		{X: -1, Y: -1},
		{X: 1, Y: 1},
	}
	blockCoords = []lostinspace.BlockCoord{
		{X: 0, Y: 1},
		{X: 15, Y: 15},
		{X: 0, Y: 0},
	}
)

func TestWorldCoord(t *testing.T) {
	for i, worldCoord := range worldCoords {
		chunkCoord, blockCoord := worldCoord.Parse()

		if chunkCoord != chunkCoords[i] {
			t.Errorf("ChunkCoord: %v != %v\n", chunkCoord, chunkCoords[i])
			t.Failed()
		}
		if blockCoord != blockCoords[i] {
			t.Errorf("BlockCoord: %v != %v\n", blockCoord, blockCoords[i])
			t.Failed()
		}

		t.Logf("%4v -> %4v, %4v\n", worldCoord, chunkCoord, blockCoord)
	}
}
