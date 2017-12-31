package lostinspace

// Terrian is set of chunks.
type Terrain struct {
	Sectors map[WorldSectorCoord]*Sector
}

func NewTerrain() *Terrain {
	terrain := &Terrain{
		Sectors: make(map[WorldSectorCoord]*Sector),
	}

	return terrain
}

func (terrain *Terrain) SetSector(sector *Sector) {
	terrain.Sectors[sector.coord] = sector
}

func (terrain *Terrain) GetSector(coord WorldSectorCoord) *Sector {
	return terrain.Sectors[coord]
}

func (terrain *Terrain) DeleteSector(coord WorldSectorCoord) {
	delete(terrain.Sectors, coord)
}

func (terrain *Terrain) SetChunk(coord WorldChunkCoord, chunk *Chunk) {
	sectorCoord, chunkCoord := coord.Parse()

	sector, exist := terrain.Sectors[sectorCoord]
	if !exist {
		// do something
		return
	}

	chunk.coord = chunkCoord
	sector.Set(chunk)
}

func (terrain *Terrain) GetChunk(coord WorldChunkCoord) *Chunk {
	sectorCoord, chunkCoord := coord.Parse()

	sector, exist := terrain.Sectors[sectorCoord]
	if !exist {
		return nil
	}

	return sector.At(chunkCoord)
}

// Set block to given world coord.
// Block coord which a block has will be ignored by world coord.
func (terrain *Terrain) SetBlock(coord WorldBlockCoord, block *Block) {
	sectorCoord, chunkCoord, blockCoord := coord.Parse()

	sector, exist := terrain.Sectors[sectorCoord]
	if !exist {
		// do something
		return
	}

	chunk := sector.At(chunkCoord)

	block.coord = blockCoord
	chunk.Set(block)
}

// Get block at given world coord.
// It will return nil if there is no corresponding chunk.
func (terrain *Terrain) GetBlock(coord WorldBlockCoord) *Block {
	sectorCoord, chunkCoord, blockCoord := coord.Parse()

	sector, exist := terrain.Sectors[sectorCoord]
	if !exist {
		return nil
	}

	chunk := sector.At(chunkCoord)

	return chunk.At(blockCoord)
}

func (terrain *Terrain) BakeChunk(world *World, dic *BlockTypeDictionary, coord WorldChunkCoord) {
	chunk := terrain.GetChunk(coord)
	if chunk == nil {
		return
	}

	BakeBlockStorageMesh(chunk.mesh, chunk, dic)

	if chunk.body == nil {
		chunk.body = world.CreateBody(false)
		chunk.body.SetPosition(
			float64(coord.X*CHUNK_WIDTH),
			float64(coord.Y*CHUNK_HEIGHT),
		)
	}
	BakeBlockStorageBody(chunk.body, chunk, dic)
}
