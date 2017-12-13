package main

// Entity is movable object which is consist of blocks.
type Entity struct {
	blocks map[Coord]*Block

	*BlockContainerObject
}

func (entity *Entity) Set(block *Block) {
	coord := Coord{int64(block.X), int64(block.Y)}
	entity.blocks[coord] = block
}

func (entity *Entity) At(x, y int64) *Block {
	return entity.blocks[Coord{x, y}]
}

func (entity *Entity) ForEach(f func(*Block)) {
	for _, v := range entity.blocks {
		f(v)
	}
}

func (entity *Entity) bake(game *Game) {
	object := NewBlockContainerObject(game, entity, BLOCK_CONTAINER_DYNAMIC)
	//object.body.SetTransform(box2d.MakeB2Vec2(
	//	float64(0),
	//	float64(0),
	//), 0)

	entity.BlockContainerObject = object
}

func NewEntity() *Entity {
	return &Entity{
		blocks: make(map[Coord]*Block),
	}
}
