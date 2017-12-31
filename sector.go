package lostinspace

// A sector is a set of chunks.
// One sector will be saved into one file.
type Sector struct {
	Chunks [SECTOR_WIDTH * SECTOR_HEIGHT]*Chunk

	coord WorldSectorCoord
}

func NewSector(coord WorldSectorCoord) *Sector {
	sector := &Sector{
		coord: coord,
	}

	return sector
}

func (sector *Sector) Set(chunk *Chunk) {
	if !chunk.coord.Valid() {
		return
	}
	sector.Chunks[chunk.coord.X+chunk.coord.Y*SECTOR_WIDTH] = chunk
}

func (sector *Sector) At(coord ChunkCoord) *Chunk {
	if !coord.Valid() {
		return nil
	}
	return sector.Chunks[coord.X+coord.Y*SECTOR_WIDTH]
}

// Deallocate all chunks.
// vao, vbo, ebo and b2body.
func (sector *Sector) Destroy() {
	for _, chunk := range sector.Chunks {
		chunk.Destroy()
	}
}
