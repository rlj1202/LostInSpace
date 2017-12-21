package lostinspace

const (
	CHUNK_WIDTH  = 16
	CHUNK_HEIGHT = 16
)

type Coord interface {
}

// WorldCoord represents absolute coordinate of block.
type WorldCoord struct {
	X, Y int64
}

// ChunkCoord represents coordinate of chunk in world.
type ChunkCoord struct {
	X, Y int64
}

// BlockCoord represents coordinate of block relative to chunk.
type BlockCoord struct {
	X, Y uint8
}

// Make WorldCoord struct object.
func MakeWorldCoord(x, y int64) WorldCoord {
	return WorldCoord{X: x, Y: y}
}

// Combind chunk coords and block coords to world coords.
func CombindWorldCoord(chunkCoord ChunkCoord, blockCoord BlockCoord) WorldCoord {
	return WorldCoord{
		X: int64(blockCoord.X) + chunkCoord.X*CHUNK_WIDTH,
		Y: int64(blockCoord.Y) + chunkCoord.Y*CHUNK_HEIGHT,
	}
}

// Get chunk coordinate and block coordinate from world coordinate.
// https://play.golang.org/p/iwLODDMvk5
func (coord *WorldCoord) Parse() (ChunkCoord, BlockCoord) {
	chunkCoord := ChunkCoord{}
	blockCoord := BlockCoord{}

	chunkCoord.X = coord.X
	chunkCoord.Y = coord.Y
	if chunkCoord.X < 0 {
		chunkCoord.X++
	}
	if chunkCoord.Y < 0 {
		chunkCoord.Y++
	}
	chunkCoord.X /= CHUNK_WIDTH
	chunkCoord.Y /= CHUNK_HEIGHT
	if coord.X < 0 {
		chunkCoord.X--
	}
	if coord.Y < 0 {
		chunkCoord.Y--
	}

	blockX := coord.X % CHUNK_WIDTH
	blockY := coord.Y % CHUNK_HEIGHT
	if blockX < 0 {
		blockX += CHUNK_WIDTH
	}
	if blockY < 0 {
		blockY += CHUNK_HEIGHT
	}
	blockCoord.X = uint8(blockX)
	blockCoord.Y = uint8(blockY)

	return chunkCoord, blockCoord
}
