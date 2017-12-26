package lostinspace

type BlockEntity struct {
	blocks map[BlockCoord]*Block

	*Mesh
	*Body
}

func NewBlockEntity(world *World) *BlockEntity {
	entity := new(BlockEntity)
	entity.blocks = make(map[BlockCoord]*Block)
	entity.Body = world.CreateBody(true)

	return entity
}

func (entity *BlockEntity) Set(block *Block) {
	entity.blocks[block.coord] = block
}

func (entity *BlockEntity) At(coord BlockCoord) *Block {
	return entity.blocks[coord]
}

func (entity *BlockEntity) ForEach(f func(*Block)) {
	for _, block := range entity.blocks {
		f(block)
	}
}

func (entity *BlockEntity) Bake(world *World, dic *BlockTypeDictionary) {
	positions, _, coords, indices := BakeBlockStorageMesh(entity, dic)

	if entity.Mesh == nil {
		entity.Mesh = NewMesh(positions, nil, coords, indices)
	} else {
		entity.Mesh.Set(positions, nil, coords, indices)
	}

	BakeBlockStorageBody(entity.Body, entity, dic)
}

func (entity *BlockEntity) Destroy() {
	entity.blocks = nil

	entity.Mesh.Destroy()
	entity.Body.Destroy()
	entity.Mesh = nil
	entity.Body = nil
}
