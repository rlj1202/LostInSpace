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
	entity.Mesh = NewMesh(nil, nil, nil, nil)

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
	BakeBlockStorageMesh(entity.Mesh, entity, dic)
	entity.Mesh.Bake()

	entity.Body.Clear()
	BakeBlockStorageBody(entity.Body, entity, dic)
	entity.Body.Bake()
}

func (entity *BlockEntity) Destroy() {
	entity.blocks = nil

	entity.Mesh.Destroy()
	entity.Body.Destroy()
	entity.Mesh = nil
	entity.Body = nil
}
