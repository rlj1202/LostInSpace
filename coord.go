package lostinspace

import (
	"fmt"
	"math"
)

const (
	SECTOR_WIDTH  = 16
	SECTOR_HEIGHT = 16

	CHUNK_WIDTH  = 16
	CHUNK_HEIGHT = 16
)

// SectorCoord represents coordinates of sector in world.
type WorldSectorCoord struct {
	X, Y int64
}

// Absolute chunk coordinates.
type WorldChunkCoord struct {
	X, Y int64
}

// WorldCoord represents absolute coordinates of block.
type WorldBlockCoord struct {
	X, Y int64
}

// ChunkCoord represents coordinates of chunk relative to sector.
type ChunkCoord struct {
	X, Y uint8
}

// BlockCoord represents coordinates of block relative to chunk.
type BlockCoord struct {
	X, Y uint8
}

func CombineWorldChunkCoord(sectorCoord WorldSectorCoord, chunkCoord ChunkCoord) WorldChunkCoord {
	return WorldChunkCoord{
		X: int64(chunkCoord.X) + sectorCoord.X*SECTOR_WIDTH,
		Y: int64(chunkCoord.Y) + sectorCoord.Y*SECTOR_HEIGHT,
	}
}

// Combind chunk coords and block coords to world coords.
func CombineWorldBlockCoord(sectorCoord WorldSectorCoord, chunkCoord ChunkCoord, blockCoord BlockCoord) WorldBlockCoord {
	return WorldBlockCoord{
		X: int64(blockCoord.X+chunkCoord.X*CHUNK_WIDTH) + sectorCoord.X*SECTOR_WIDTH*CHUNK_WIDTH,
		Y: int64(blockCoord.Y+chunkCoord.Y*CHUNK_HEIGHT) + sectorCoord.Y*SECTOR_HEIGHT*CHUNK_HEIGHT,
	}
}

func (coord *WorldChunkCoord) Parse() (WorldSectorCoord, ChunkCoord) {
	sectorCoord := WorldSectorCoord{
		X: int64(math.Floor(float64(coord.X) / float64(SECTOR_WIDTH))),
		Y: int64(math.Floor(float64(coord.Y) / float64(SECTOR_HEIGHT))),
	}
	sectorXOff := coord.X % SECTOR_WIDTH
	sectorYOff := coord.Y % SECTOR_HEIGHT
	if sectorXOff < 0 {
		sectorXOff += SECTOR_WIDTH
	}
	if sectorYOff < 0 {
		sectorYOff += SECTOR_HEIGHT
	}

	chunkCoord := ChunkCoord{
		X: uint8(sectorXOff),
		Y: uint8(sectorYOff),
	}

	return sectorCoord, chunkCoord
}

// Get chunk coordinate and block coordinate from world coordinate.
func (coord *WorldBlockCoord) Parse() (WorldSectorCoord, ChunkCoord, BlockCoord) {
	sectorCoord := WorldSectorCoord{
		X: int64(math.Floor(float64(coord.X) / float64(SECTOR_WIDTH*CHUNK_WIDTH))),
		Y: int64(math.Floor(float64(coord.Y) / float64(SECTOR_HEIGHT*CHUNK_HEIGHT))),
	}
	sectorXOff := coord.X % (SECTOR_WIDTH * CHUNK_WIDTH)
	sectorYOff := coord.Y % (SECTOR_HEIGHT * CHUNK_HEIGHT)
	if sectorXOff < 0 {
		sectorXOff += SECTOR_WIDTH * CHUNK_WIDTH
	}
	if sectorYOff < 0 {
		sectorYOff += SECTOR_HEIGHT * CHUNK_HEIGHT
	}

	chunkCoord := ChunkCoord{
		X: uint8(math.Floor(float64(sectorXOff) / float64(CHUNK_WIDTH))),
		Y: uint8(math.Floor(float64(sectorYOff) / float64(CHUNK_HEIGHT))),
	}
	chunkXOff := sectorXOff % CHUNK_WIDTH
	chunkYOff := sectorYOff % CHUNK_HEIGHT
	if chunkXOff < 0 {
		chunkXOff += CHUNK_WIDTH
	}
	if chunkYOff < 0 {
		chunkYOff += CHUNK_HEIGHT
	}

	blockCoord := BlockCoord{
		X: uint8(chunkXOff),
		Y: uint8(chunkYOff),
	}

	return sectorCoord, chunkCoord, blockCoord
}

func (coord WorldSectorCoord) String() string {
	return fmt.Sprintf("WorldSectorCoord{X: %d, Y: %d}", coord.X, coord.Y)
}

func (coord WorldChunkCoord) String() string {
	return fmt.Sprintf("WorldChunkCoord{X: %d, Y: %d}", coord.X, coord.Y)
}

func (coord WorldBlockCoord) String() string {
	return fmt.Sprintf("WorldBlockCoord{X: %d, Y: %d}", coord.X, coord.Y)
}

func (coord ChunkCoord) String() string {
	return fmt.Sprintf("ChunkCoord{X: %d, Y: %d}", coord.X, coord.Y)
}

func (coord BlockCoord) String() string {
	return fmt.Sprintf("BlockCoord{X: %d, Y: %d}", coord.X, coord.Y)
}

func (coord ChunkCoord) Valid() bool {
	return 0 <= coord.X && coord.X < SECTOR_WIDTH && 0 <= coord.Y && coord.Y < SECTOR_HEIGHT
}

func (coord BlockCoord) Valid() bool {
	return 0 <= coord.X && coord.X < CHUNK_WIDTH && 0 <= coord.Y && coord.Y < CHUNK_HEIGHT
}

func (coord WorldSectorCoord) Left() WorldSectorCoord {
	return WorldSectorCoord{coord.X - 1, coord.Y}
}

func (coord WorldSectorCoord) Right() WorldSectorCoord {
	return WorldSectorCoord{coord.X + 1, coord.Y}
}

func (coord WorldSectorCoord) Up() WorldSectorCoord {
	return WorldSectorCoord{coord.X, coord.Y + 1}
}

func (coord WorldSectorCoord) Down() WorldSectorCoord {
	return WorldSectorCoord{coord.X, coord.Y - 1}
}
