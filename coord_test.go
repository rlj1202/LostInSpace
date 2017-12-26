package lostinspace_test

import (
	"github.com/rlj1202/LostInSpace"
	"testing"
)

var (
	worldCoords = []lostinspace.WorldBlockCoord{
		{X: 0, Y: 1},
		{X: -1, Y: -1},
		{X: 16, Y: 16},
		{X: 256, Y: 256},
	}
	sectorCoords = []lostinspace.WorldSectorCoord{
		{X: 0, Y: 0},
		{X: -1, Y: -1},
		{X: 0, Y: 0},
		{X: 1, Y: 1},
	}
	chunkCoords = []lostinspace.ChunkCoord{
		{X: 0, Y: 0},
		{X: 15, Y: 15},
		{X: 1, Y: 1},
		{X: 0, Y: 0},
	}
	blockCoords = []lostinspace.BlockCoord{
		{X: 0, Y: 1},
		{X: 15, Y: 15},
		{X: 0, Y: 0},
		{X: 0, Y: 0},
	}
)

func TestWorldCoord(t *testing.T) {
	for i, worldCoord := range worldCoords {
		sectorCoord, chunkCoord, blockCoord := worldCoord.Parse()

		if sectorCoord != sectorCoords[i] {
			t.Errorf("SectorCoord: %v != %v\n", sectorCoord, sectorCoords[i])
			t.Failed()
		}
		if chunkCoord != chunkCoords[i] {
			t.Errorf("ChunkCoord: %v != %v\n", chunkCoord, chunkCoords[i])
			t.Failed()
		}
		if blockCoord != blockCoords[i] {
			t.Errorf("BlockCoord: %v != %v\n", blockCoord, blockCoords[i])
			t.Failed()
		}

		t.Logf("%4v -> %4v, %4v, %4v\n", worldCoord, sectorCoord, chunkCoord, blockCoord)

		newWorldCoord := lostinspace.CombineWorldBlockCoord(sectorCoord, chunkCoord, blockCoord)

		if worldCoord != newWorldCoord {
			t.Errorf("NewWorldCoord: %v != %v\n", worldCoord, newWorldCoord)
		}

		t.Logf("%4v, %4v, %4v -> %4v\n", sectorCoord, chunkCoord, blockCoord, newWorldCoord)
	}

	test := lostinspace.WorldChunkCoord{-1, -1}
	a, b := test.Parse()
	t.Log(a, b)
}
